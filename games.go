package main

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/msg"
	"google.golang.org/protobuf/proto"
)

const (
	mapIconBase    = "https://raw.githubusercontent.com/MurkyYT/cs2-map-icons/main/images/"
	maxGameScan    = 40
	mapParseFrames = 64
	gameWorkers    = 4
)

var (
	errNoGameFolder = errors.New("set the CS2 game folder first")
	errBadGameID    = errors.New("invalid demo path")

	autoDemoDateRe = regexp.MustCompile(`(?:^|[-_])((?:19|20)\d{2})[-_](\d{2})[-_](\d{2})(?:[-_](\d{2})[-_](\d{2})[-_](\d{2}))?`)
	unixNameRe     = regexp.MustCompile(`^(\d{10})(?:[_-]|\.dem)`)
	mapInNameRe    = regexp.MustCompile(`(?:^|[^a-z0-9])((?:de|cs|ar|gd)_[a-z0-9]+)`)
)

var skipDemoDirs = map[string]bool{
	"maps": true, "materials": true, "models": true, "panorama": true,
	"sound": true, "resource": true, "scripts": true, "cfg": true,
	"bin": true, "shaders": true, "particles": true, "characters": true,
	"weapons": true, "ui": true, "addons": true, "media": true,
	"scenes": true, "expressions": true, "workshop": true,
}

var mapTitles = map[string]string{
	"de_ancient": "Ancient", "de_anubis": "Anubis", "de_cache": "Cache",
	"de_dust": "Dust", "de_dust2": "Dust II", "de_inferno": "Inferno",
	"de_mirage": "Mirage", "de_nuke": "Nuke", "de_overpass": "Overpass",
	"de_train": "Train", "de_vertigo": "Vertigo", "cs_office": "Office",
	"cs_italy": "Italy", "cs_agency": "Agency", "de_mills": "Mills",
	"de_thera": "Thera", "de_grail": "Grail", "de_basalt": "Basalt",
	"de_edin": "Edin", "de_dogtown": "Dogtown",
}

type demoFile struct {
	Path    string
	Rel     string
	Name    string
	ModTime time.Time
}

type pastGame struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Map      string `json:"map"`
	MapTitle string `json:"mapTitle"`
	IconURL  string `json:"iconUrl"`
	DateUnix int64  `json:"dateUnix"`
	Error    string `json:"error,omitempty"`
}

type gamesResponse struct {
	Games     []pastGame `json:"games"`
	Truncated bool       `json:"truncated"`
}

type gameScoreResponse struct {
	CT int `json:"ct"`
	T  int `json:"t"`
}

func (a *webApp) handleGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mu.Lock()
	folder := a.cfg.GameFolder
	a.mu.Unlock()
	if folder == "" {
		http.Error(w, errNoGameFolder.Error(), http.StatusBadRequest)
		return
	}
	files, err := listDemoFiles(folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	truncated := false
	if len(files) > maxGameScan {
		files = files[:maxGameScan]
		truncated = true
	}
	writeJSON(w, gamesResponse{
		Games:     buildPastGames(files),
		Truncated: truncated,
	})
}

func (a *webApp) handleGameScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mu.Lock()
	folder := a.cfg.GameFolder
	a.mu.Unlock()
	if folder == "" {
		http.Error(w, errNoGameFolder.Error(), http.StatusBadRequest)
		return
	}
	full, err := demoPathInFolder(folder, r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ct, tSide, err := parseDemoScore(full)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, gameScoreResponse{CT: ct, T: tSide})
}

func (a *webApp) handleGameOpen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mu.Lock()
	folder := a.cfg.GameFolder
	a.mu.Unlock()
	if folder == "" {
		http.Error(w, errNoGameFolder.Error(), http.StatusBadRequest)
		return
	}
	full, err := demoPathInFolder(folder, r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := a.loadFromPath(full, filepath.Base(full)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.mu.Lock()
	a.copied = true
	a.copiedTo = full
	a.copyErr = ""
	a.mu.Unlock()
	writeJSON(w, a.snapshot())
}

func listDemoFiles(folder string) ([]demoFile, error) {
	var files []demoFile
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if path != folder && skipDemoDirs[strings.ToLower(d.Name())] {
				return filepath.SkipDir
			}
			rel, relErr := filepath.Rel(folder, path)
			if relErr == nil && rel != "." && strings.Count(rel, string(os.PathSeparator)) >= 1 {
				return filepath.SkipDir
			}
			return nil
		}
		if !isDemoFilename(d.Name()) {
			return nil
		}
		rel, relErr := filepath.Rel(folder, path)
		if relErr != nil {
			return nil
		}
		info, statErr := d.Info()
		if statErr != nil {
			return nil
		}
		files = append(files, demoFile{
			Path:    path,
			Rel:     filepath.ToSlash(rel),
			Name:    d.Name(),
			ModTime: info.ModTime(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].ModTime.Equal(files[j].ModTime) {
			return files[i].Name > files[j].Name
		}
		return files[i].ModTime.After(files[j].ModTime)
	})
	return files, nil
}

func buildPastGames(files []demoFile) []pastGame {
	out := make([]pastGame, len(files))
	var wg sync.WaitGroup
	sem := make(chan struct{}, gameWorkers)
	for i, f := range files {
		wg.Add(1)
		go func(i int, f demoFile) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			out[i] = buildPastGame(f)
		}(i, f)
	}
	wg.Wait()
	return out
}

func buildPastGame(f demoFile) pastGame {
	g := pastGame{
		ID:   f.Rel,
		Name: f.Name,
	}
	if t := demoDate(f); !t.IsZero() {
		g.DateUnix = t.Unix()
	}
	mapName, err := parseDemoMap(f.Path)
	if err != nil {
		g.Error = err.Error()
	}
	if mapName == "" {
		mapName = mapNameFromFilename(f.Name)
	}
	g.Map = mapName
	g.MapTitle = mapTitle(mapName)
	g.IconURL = mapIconURL(mapName)
	return g
}

func demoDate(f demoFile) time.Time {
	if t, ok := parseDateFromInfoFile(f.Path); ok {
		return t
	}
	if t, ok := parseDateFromName(f.Name); ok {
		return t
	}
	return f.ModTime
}

func parseDateFromName(name string) (time.Time, bool) {
	if m := unixNameRe.FindStringSubmatch(name); m != nil {
		sec, err := strconv.ParseInt(m[1], 10, 64)
		if err == nil && sec >= 1420070400 && sec <= time.Now().Add(24*time.Hour).Unix() {
			return time.Unix(sec, 0), true
		}
	}
	m := autoDemoDateRe.FindStringSubmatch(name)
	if m == nil {
		return time.Time{}, false
	}
	layout := "2006-01-02"
	value := m[1] + "-" + m[2] + "-" + m[3]
	if m[4] != "" {
		layout = "2006-01-02 15-04-05"
		value += " " + m[4] + "-" + m[5] + "-" + m[6]
	}
	t, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func parseDateFromInfoFile(demoPath string) (time.Time, bool) {
	for _, p := range infoFileCandidates(demoPath) {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var info msg.CDataGCCStrike15V2_MatchInfo
		if err := proto.Unmarshal(data, &info); err != nil {
			continue
		}
		if info.GetMatchtime() == 0 {
			continue
		}
		return time.Unix(int64(info.GetMatchtime()), 0), true
	}
	return time.Time{}, false
}

func infoFileCandidates(demoPath string) []string {
	out := []string{demoPath + ".info"}
	if strings.HasSuffix(strings.ToLower(demoPath), ".zst") {
		out = append(out, strings.TrimSuffix(demoPath, filepath.Ext(demoPath))+".info")
	}
	return out
}

func mapNameFromFilename(name string) string {
	lower := strings.ToLower(name)
	lower = strings.TrimSuffix(lower, ".zst")
	lower = strings.TrimSuffix(lower, ".dem")
	m := mapInNameRe.FindStringSubmatch(lower)
	if m == nil {
		return ""
	}
	return m[1]
}

func mapTitle(name string) string {
	if name == "" {
		return ""
	}
	if t, ok := mapTitles[name]; ok {
		return t
	}
	s := name
	for _, p := range []string{"de_", "cs_", "ar_", "gd_"} {
		s = strings.TrimPrefix(s, p)
	}
	parts := strings.Fields(strings.ReplaceAll(s, "_", " "))
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

func mapIconURL(mapName string) string {
	if mapName == "" {
		return ""
	}
	return mapIconBase + mapName + ".png"
}

func cleanMapName(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	raw = strings.ReplaceAll(raw, `\`, "/")
	if i := strings.LastIndex(raw, "/"); i >= 0 {
		raw = raw[i+1:]
	}
	return raw
}

func parseDemoMap(rawPath string) (mapName string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New("could not read demo header")
		}
	}()
	demoPath, cleanup, err := resolveDemoPath(rawPath)
	if err != nil {
		return "", err
	}
	if cleanup != nil {
		defer cleanup()
	}
	f, err := os.Open(demoPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()

	var (
		got string
		mu  sync.Mutex
	)
	setMap := func(name string) {
		mu.Lock()
		if got == "" {
			got = name
		}
		mu.Unlock()
	}
	p.RegisterNetMessageHandler(func(m *msg.CSVCMsg_ServerInfo) {
		setMap(m.GetMapName())
	})
	p.RegisterNetMessageHandler(func(m *msg.CNETMsg_SignonState) {
		setMap(m.GetMapName())
	})
	for i := 0; i < mapParseFrames; i++ {
		mu.Lock()
		have := got
		mu.Unlock()
		if have != "" {
			break
		}
		more, frameErr := p.ParseNextFrame()
		if frameErr != nil || !more {
			break
		}
	}
	mu.Lock()
	mapName = got
	mu.Unlock()
	return cleanMapName(mapName), nil
}

func parseDemoScore(rawPath string) (ct, t int, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New("could not read demo score")
		}
	}()
	demoPath, cleanup, err := resolveDemoPath(rawPath)
	if err != nil {
		return 0, 0, err
	}
	if cleanup != nil {
		defer cleanup()
	}
	f, err := os.Open(demoPath)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()
	if err := p.ParseToEnd(); err != nil && !errors.Is(err, demoinfocs.ErrUnexpectedEndOfDemo) {
		return 0, 0, err
	}
	gs := p.GameState()
	return gs.TeamCounterTerrorists().Score(), gs.TeamTerrorists().Score(), nil
}

func demoPathInFolder(folder, id string) (string, error) {
	id = strings.TrimSpace(id)
	if id == "" || strings.Contains(id, "..") {
		return "", errBadGameID
	}
	full := filepath.Join(folder, filepath.FromSlash(id))
	absFolder, err := filepath.Abs(folder)
	if err != nil {
		return "", err
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	sep := string(os.PathSeparator)
	if absFull != absFolder && !strings.HasPrefix(absFull, absFolder+sep) {
		return "", errBadGameID
	}
	if !isDemoFilename(filepath.Base(absFull)) {
		return "", errBadGameID
	}
	info, err := os.Stat(absFull)
	if err != nil || info.IsDir() {
		return "", errBadGameID
	}
	return absFull, nil
}
