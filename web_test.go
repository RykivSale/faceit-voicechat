package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	app.handleIndex(w, r)
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "/api/state") || !strings.Contains(body, "faceit-voicechat") {
		t.Fatal("index html missing expected content")
	}
}

func TestHandleStateEmpty(t *testing.T) {
	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	w := httptest.NewRecorder()
	app.handleState(w, r)
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	var got stateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Demo.Loaded {
		t.Fatal("expected no demo")
	}
	if got.DetectedFolders == nil {
		t.Fatal("detectedFolders should be [] not null")
	}
	if got.Config.Keys[0] != "F5" {
		t.Fatalf("keys=%v", got.Config.Keys)
	}
}

func TestHandleStateWithDemo(t *testing.T) {
	app := &webApp{
		cfg:      defaultConfig(),
		shutdown: make(chan struct{}),
		parsed:   true,
		demoName: "faceit.dem",
		demoPath: "/tmp/faceit.dem",
		ctMask:   7,
		tMask:    24,
	}
	r := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	w := httptest.NewRecorder()
	app.handleState(w, r)
	var got stateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if !got.Demo.Loaded || got.Demo.Bind == "" || got.Demo.Playdemo != "playdemo faceit" {
		t.Fatalf("unexpected demo view: %+v", got.Demo)
	}
	if !strings.Contains(got.Demo.Bind, `bind "F5"`) {
		t.Fatalf("bind=%s", got.Demo.Bind)
	}
}

func TestHandleConfigRejectsMissingFolder(t *testing.T) {
	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	body := `{"gameFolder":"/definitely/not/a/real/cs2/path-xyz"}`
	r := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
	w := httptest.NewRecorder()
	app.handleConfig(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
}

func TestHandleConfigSavesKeys(t *testing.T) {
	dir := t.TempDir()
	old := configPathFn
	configPathFn = func() (string, error) {
		return filepath.Join(dir, "config.json"), nil
	}
	t.Cleanup(func() { configPathFn = old })

	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	body := `{"keys":["kp_ins","kp_del","kp_end"]}`
	r := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
	w := httptest.NewRecorder()
	app.handleConfig(w, r)
	if w.Code != 200 {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
	var got stateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Config.Keys[0] != "KP_INS" || got.Config.Keys[2] != "KP_END" {
		t.Fatalf("keys=%v", got.Config.Keys)
	}
}

func TestHandleDemoRejectsNonDemo(t *testing.T) {
	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("demo", "notes.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(fw, "hello")
	mw.Close()

	r := httptest.NewRequest(http.MethodPost, "/api/demo", &buf)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	app.handleDemo(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.dem")
	dst := filepath.Join(dir, "b.dem")
	if err := os.WriteFile(src, []byte("demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "demo" {
		t.Fatalf("got %q", got)
	}
}

func TestHandleQuit(t *testing.T) {
	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodPost, "/api/quit", nil)
	w := httptest.NewRecorder()
	app.handleQuit(w, r)
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
}
