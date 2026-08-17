package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed web
var webFS embed.FS

const maxDemoUpload = 1 << 30 // 1 GiB

var errBadFolder = errors.New("that folder does not exist")

type webApp struct {
	mu       sync.Mutex
	cfg      config
	demoName string
	demoPath string
	ctMask   uint32
	tMask    uint32
	parsed   bool
	parseErr string
	copied   bool
	copiedTo string
	copyErr  string
	cleanup  func()

	shutdown     chan struct{}
	shutdownOnce sync.Once
}

type stateResponse struct {
	Config              config   `json:"config"`
	Demo                demoView `json:"demo"`
	DetectedFolders     []string `json:"detectedFolders"`
	LooksLikeGameFolder bool     `json:"looksLikeGameFolder"`
}

type demoView struct {
	Loaded   bool   `json:"loaded"`
	Name     string `json:"name"`
	Error    string `json:"error"`
	Bind     string `json:"bind"`
	Playdemo string `json:"playdemo"`
	Copied   bool   `json:"copied"`
	CopiedTo string `json:"copiedTo"`
	CopyErr  string `json:"copyError"`
}

func runWeb(rawPath string) {
	app := &webApp{
		cfg:      loadConfig(),
		shutdown: make(chan struct{}),
	}

	if rawPath != "" {
		fmt.Println("Reading demo...")
		if err := app.loadFromPath(rawPath, filepath.Base(rawPath)); err != nil {
			fmt.Println("Failed to open demo:", err)
			app.mu.Lock()
			app.parseErr = err.Error()
			app.demoName = filepath.Base(rawPath)
			app.mu.Unlock()
		} else {
			app.mu.Lock()
			app.tryCopyLocked()
			name := app.demoName
			app.mu.Unlock()
			fmt.Println("Demo ready:", name)
		}
	}

	ln, err := net.Listen("tcp", "127.0.0.1:8765")
	if err != nil {
		ln, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			fmt.Println("Failed to start local server:", err)
			waitForExit()
			return
		}
	}

	mux := http.NewServeMux()
	webRoot, err := fs.Sub(webFS, "web")
	if err != nil {
		fmt.Println("Failed to load web assets:", err)
		waitForExit()
		return
	}
	mux.Handle("/vendor/", http.FileServer(http.FS(webRoot)))
	mux.HandleFunc("/", app.handleIndex)
	mux.HandleFunc("/api/state", app.handleState)
	mux.HandleFunc("/api/demo", app.handleDemo)
	mux.HandleFunc("/api/config", app.handleConfig)
	mux.HandleFunc("/api/copy", app.handleCopy)
	mux.HandleFunc("/api/detect", app.handleDetect)
	mux.HandleFunc("/api/find", app.handleFind)
	mux.HandleFunc("/api/browse", app.handleBrowse)
	mux.HandleFunc("/api/games", app.handleGames)
	mux.HandleFunc("/api/games/score", app.handleGameScore)
	mux.HandleFunc("/api/games/open", app.handleGameOpen)
	mux.HandleFunc("/api/update", app.handleUpdate)
	mux.HandleFunc("/api/quit", app.handleQuit)

	srv := &http.Server{Handler: mux}
	url := "http://" + ln.Addr().String() + "/"

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
			app.requestShutdown()
		}
	}()

	fmt.Println()
	fmt.Println("Local page:", url)
	fmt.Println("Leave this window open while you use it. Press Ctrl+C to quit.")
	fmt.Println()
	openBrowser(url)

	if app.cfg.CheckUpdates {
		go func() {
			check := checkUpdatesFn()
			if check.Newer {
				printUpdateBanner(check)
			}
		}()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	select {
	case <-sig:
	case <-app.shutdown:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	app.clearDemo()
}

func (a *webApp) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(webFS, "web/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(data)
}

func (a *webApp) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, a.snapshot())
}

func (a *webApp) handleDetect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, a.snapshot())
}

func (a *webApp) handleFind(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	found := detectCS2FoldersFn()
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(found) > 0 {
		if err := a.applyGameFolderLocked(found[0]); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	writeJSON(w, a.snapshotLocked())
}

func (a *webApp) handleBrowse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path, err := pickFolderFn()
	if err != nil {
		if errors.Is(err, errPickCancelled) {
			writeJSON(w, a.snapshot())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.applyGameFolderLocked(path); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, a.snapshotLocked())
}

func (a *webApp) applyGameFolderLocked(folder string) error {
	folder = strings.Trim(strings.TrimSpace(folder), `"`)
	if folder != "" {
		folder = normalizeGameFolder(folder)
		info, err := os.Stat(folder)
		if err != nil || !info.IsDir() {
			return errBadFolder
		}
	}
	a.cfg.GameFolder = folder
	a.copied = false
	a.copiedTo = ""
	a.copyErr = ""
	if err := saveConfig(a.cfg); err != nil {
		return err
	}
	if a.parsed {
		a.tryCopyLocked()
	}
	return nil
}

func (a *webApp) handleDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxDemoUpload)
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, "could not read file (max 1 GB): "+err.Error(), http.StatusBadRequest)
		return
	}
	file, hdr, err := r.FormFile("demo")
	if err != nil {
		http.Error(w, "choose a .dem or .dem.zst file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !isDemoFilename(hdr.Filename) {
		http.Error(w, "file must be .dem or .dem.zst", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	folder := a.cfg.GameFolder
	a.mu.Unlock()
	if folder == "" {
		http.Error(w, "set the CS2 game folder first", http.StatusBadRequest)
		return
	}

	tmpDir, err := os.MkdirTemp("", "faceit-voicechat-upload-*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rawPath := filepath.Join(tmpDir, filepath.Base(hdr.Filename))
	out, err := os.Create(rawPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		os.RemoveAll(tmpDir)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out.Close()

	if err := a.loadFromPath(rawPath, hdr.Filename); err != nil {
		os.RemoveAll(tmpDir)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	prev := a.cleanup
	a.cleanup = func() {
		if prev != nil {
			prev()
		}
		os.RemoveAll(tmpDir)
	}
	a.tryCopyLocked()
	a.mu.Unlock()

	writeJSON(w, a.snapshot())
}

func (a *webApp) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		GameFolder *string  `json:"gameFolder"`
		Keys       []string `json:"keys"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if body.GameFolder != nil {
		if err := a.applyGameFolderLocked(*body.GameFolder); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if len(body.Keys) > 0 {
		if len(body.Keys) != 3 {
			http.Error(w, "please enter exactly 3 keys", http.StatusBadRequest)
			return
		}
		keys := make([]string, 0, 3)
		for _, k := range body.Keys {
			key := strings.ToUpper(strings.TrimSpace(k))
			if key == "" {
				http.Error(w, "keys cannot be empty", http.StatusBadRequest)
				return
			}
			keys = append(keys, key)
		}
		a.cfg.Keys = keys
	}

	if err := saveConfig(a.cfg); err != nil {
		http.Error(w, "failed to save settings: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, a.snapshotLocked())
}

func (a *webApp) handleCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.parsed {
		http.Error(w, "open a demo first", http.StatusBadRequest)
		return
	}
	a.tryCopyLocked()
	if a.copyErr != "" && !a.copied {
		http.Error(w, a.copyErr, http.StatusBadRequest)
		return
	}
	writeJSON(w, a.snapshotLocked())
}

func (a *webApp) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mu.Lock()
	enabled := a.cfg.CheckUpdates
	a.mu.Unlock()
	if !enabled {
		writeJSON(w, updateView{Current: appVersion, URL: githubLatestURL})
		return
	}
	writeJSON(w, toUpdateView(checkUpdatesFn()))
}

func (a *webApp) handleQuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true}`))
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.requestShutdown()
	}()
}

func (a *webApp) loadFromPath(rawPath, displayName string) error {
	demoPath, cleanup, err := resolveDemoPath(rawPath)
	if err != nil {
		return err
	}
	ctMask, tMask, err := parseDemo(demoPath)
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.clearDemoLocked()
	a.demoPath = demoPath
	a.demoName = displayName
	a.ctMask = ctMask
	a.tMask = tMask
	a.parsed = true
	a.parseErr = ""
	a.copied = false
	a.copiedTo = ""
	a.copyErr = ""
	a.cleanup = cleanup
	return nil
}

func (a *webApp) tryCopyLocked() {
	if !a.parsed || a.cfg.GameFolder == "" {
		return
	}
	dst, err := copyDemoToGameFolder(a.cfg, a.demoPath)
	if err != nil {
		a.copied = false
		a.copiedTo = ""
		a.copyErr = err.Error()
		return
	}
	a.copied = true
	a.copiedTo = dst
	a.copyErr = ""
}

func (a *webApp) clearDemo() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.clearDemoLocked()
}

func (a *webApp) clearDemoLocked() {
	if a.cleanup != nil {
		a.cleanup()
		a.cleanup = nil
	}
	a.demoName = ""
	a.demoPath = ""
	a.ctMask = 0
	a.tMask = 0
	a.parsed = false
	a.parseErr = ""
	a.copied = false
	a.copiedTo = ""
	a.copyErr = ""
}

func (a *webApp) snapshot() stateResponse {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.snapshotLocked()
}

func (a *webApp) snapshotLocked() stateResponse {
	view := demoView{
		Loaded:   a.parsed,
		Name:     a.demoName,
		Error:    a.parseErr,
		Copied:   a.copied,
		CopiedTo: a.copiedTo,
		CopyErr:  a.copyErr,
	}
	if a.parsed {
		view.Bind = bindCommand(a.cfg, a.ctMask, a.tMask)
		src := a.demoPath
		if a.copied && a.copiedTo != "" {
			src = a.copiedTo
		}
		view.Playdemo = playdemoCommand(a.cfg.GameFolder, src)
	}
	looksLike := a.cfg.GameFolder == "" || looksLikeCSGOFolder(a.cfg.GameFolder)
	return stateResponse{
		Config:              a.cfg,
		Demo:                view,
		DetectedFolders:     detectCS2FoldersFn(),
		LooksLikeGameFolder: looksLike,
	}
}

func (a *webApp) requestShutdown() {
	a.shutdownOnce.Do(func() { close(a.shutdown) })
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
