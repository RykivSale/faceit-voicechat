package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/msg"
	"google.golang.org/protobuf/proto"
)

func TestMapIconURL(t *testing.T) {
	if got := mapIconURL("de_mirage"); got != "https://raw.githubusercontent.com/MurkyYT/cs2-map-icons/main/images/de_mirage.png" {
		t.Fatalf("got %q", got)
	}
	if mapIconURL("") != "" {
		t.Fatal("empty map should have no icon")
	}
}

func TestMapTitle(t *testing.T) {
	if got := mapTitle("de_dust2"); got != "Dust II" {
		t.Fatalf("dust2=%q", got)
	}
	if got := mapTitle("de_mirage"); got != "Mirage" {
		t.Fatalf("mirage=%q", got)
	}
	if got := mapTitle("de_some_custom"); got != "Some Custom" {
		t.Fatalf("fallback=%q", got)
	}
	if mapTitle("") != "" {
		t.Fatal("empty")
	}
}

func TestMapNameFromFilename(t *testing.T) {
	if got := mapNameFromFilename("auto-2024-01-02-03-04-05-de_mirage.dem"); got != "de_mirage" {
		t.Fatalf("got %q", got)
	}
	if got := mapNameFromFilename("faceit-de_dust2.dem.zst"); got != "de_dust2" {
		t.Fatalf("got %q", got)
	}
	if mapNameFromFilename("match730_123.dem") != "" {
		t.Fatal("expected no map")
	}
}

func TestParseDateFromName(t *testing.T) {
	got, ok := parseDateFromName("auto-2024-08-17-21-30-00-de_mirage.dem")
	if !ok {
		t.Fatal("expected date")
	}
	if got.Year() != 2024 || got.Month() != 8 || got.Day() != 17 || got.Hour() != 21 || got.Minute() != 30 {
		t.Fatalf("got %v", got)
	}
	if _, ok := parseDateFromName("faceit-abc.dem"); ok {
		t.Fatal("expected no date")
	}
	unix, ok := parseDateFromName("1692288000_team-a-team-b_de_mirage.dem")
	if !ok || unix.Unix() != 1692288000 {
		t.Fatalf("unix=%v ok=%v", unix, ok)
	}
}

func TestDemoDateFallsBackToFileTime(t *testing.T) {
	mtime := time.Date(2024, 5, 1, 12, 0, 0, 0, time.Local)
	f := demoFile{Name: "1-abc.dem", Path: "/nope/1-abc.dem", ModTime: mtime}
	got := demoDate(f)
	if !got.Equal(mtime) {
		t.Fatalf("got %v want %v", got, mtime)
	}
}

func TestParseDateFromInfoFile(t *testing.T) {
	dir := t.TempDir()
	dem := filepath.Join(dir, "match.dem")
	if err := os.WriteFile(dem, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	ts := uint32(1710518400) // 2024-03-15 16:00 UTC
	info := &msg.CDataGCCStrike15V2_MatchInfo{Matchtime: &ts}
	data, err := proto.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dem+".info", data, 0o644); err != nil {
		t.Fatal(err)
	}
	got, ok := parseDateFromInfoFile(dem)
	if !ok || got.Unix() != int64(ts) {
		t.Fatalf("got=%v ok=%v", got, ok)
	}
}

func TestCleanMapName(t *testing.T) {
	if got := cleanMapName(`workshop/123/de_cache`); got != "de_cache" {
		t.Fatalf("got %q", got)
	}
	if got := cleanMapName("DE_MIRAGE"); got != "de_mirage" {
		t.Fatalf("got %q", got)
	}
}

func TestListDemoFiles(t *testing.T) {
	root := t.TempDir()
	replays := filepath.Join(root, "replays")
	mapsDir := filepath.Join(root, "maps")
	nested := filepath.Join(replays, "old")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(mapsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite := func(path string) {
		t.Helper()
		if err := os.WriteFile(path, []byte("demo"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite(filepath.Join(root, "faceit.dem"))
	mustWrite(filepath.Join(root, "notes.txt"))
	mustWrite(filepath.Join(replays, "last.dem"))
	mustWrite(filepath.Join(nested, "hidden.dem"))
	mustWrite(filepath.Join(mapsDir, "de_mirage.dem"))

	old := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(filepath.Join(root, "faceit.dem"), old, old); err != nil {
		t.Fatal(err)
	}

	got, err := listDemoFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d want 2: %+v", len(got), got)
	}
	if got[0].Name != "last.dem" || got[0].Rel != "replays/last.dem" {
		t.Fatalf("newest=%+v", got[0])
	}
	if got[1].Name != "faceit.dem" {
		t.Fatalf("older=%+v", got[1])
	}
}

func TestDemoPathInFolder(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "a.dem")
	if err := os.WriteFile(src, []byte("demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := demoPathInFolder(root, "a.dem")
	if err != nil {
		t.Fatal(err)
	}
	if got != src {
		t.Fatalf("got %q want %q", got, src)
	}
	if _, err := demoPathInFolder(root, "../a.dem"); err == nil {
		t.Fatal("expected traversal reject")
	}
	if _, err := demoPathInFolder(root, "notes.txt"); err == nil {
		t.Fatal("expected non-demo reject")
	}
}

func TestHandleGamesRequiresFolder(t *testing.T) {
	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	w := httptest.NewRecorder()
	app.handleGames(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
}

func TestHandleGamesListsDemos(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "auto-2024-01-02-03-04-05-de_nuke.dem"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	app := &webApp{cfg: config{GameFolder: root, Keys: defaultConfig().Keys}, shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	w := httptest.NewRecorder()
	app.handleGames(w, r)
	if w.Code != 200 {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
	var got gamesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Games) != 1 {
		t.Fatalf("games=%+v", got.Games)
	}
	g := got.Games[0]
	if g.Map != "de_nuke" || g.MapTitle != "Nuke" {
		t.Fatalf("map=%q title=%q", g.Map, g.MapTitle)
	}
	if !strings.Contains(g.IconURL, "de_nuke.png") {
		t.Fatalf("icon=%q", g.IconURL)
	}
	if time.Unix(g.DateUnix, 0).Year() != 2024 {
		t.Fatalf("date=%d", g.DateUnix)
	}
}

func TestHandleGameScoreRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	app := &webApp{cfg: config{GameFolder: root, Keys: defaultConfig().Keys}, shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodGet, "/api/games/score?id=../secret.dem", nil)
	w := httptest.NewRecorder()
	app.handleGameScore(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
}

func TestHandleGameOpenRequiresFolder(t *testing.T) {
	app := &webApp{cfg: defaultConfig(), shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodPost, "/api/games/open?id=a.dem", nil)
	w := httptest.NewRecorder()
	app.handleGameOpen(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
}

func TestHandleGameOpenRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	app := &webApp{cfg: config{GameFolder: root, Keys: defaultConfig().Keys}, shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodPost, "/api/games/open?id=../secret.dem", nil)
	w := httptest.NewRecorder()
	app.handleGameOpen(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
}

func TestHandleGameOpenDoesNotCopyInvalidDemo(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "match.dem")
	if err := os.WriteFile(src, []byte("not-a-demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	app := &webApp{cfg: config{GameFolder: root, Keys: defaultConfig().Keys}, shutdown: make(chan struct{})}
	r := httptest.NewRequest(http.MethodPost, "/api/games/open?id=match.dem", nil)
	w := httptest.NewRecorder()
	app.handleGameOpen(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "not-a-demo" {
		t.Fatalf("file was rewritten: %q", got)
	}
}
