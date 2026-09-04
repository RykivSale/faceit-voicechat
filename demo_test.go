package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIdsToBitmask(t *testing.T) {
	if got := idsToBitmask(nil); got != 0 {
		t.Fatalf("empty = %d, want 0", got)
	}
	if got := idsToBitmask([]int{1}); got != 1 {
		t.Fatalf("id 1 = %d, want 1", got)
	}
	if got := idsToBitmask([]int{2}); got != 2 {
		t.Fatalf("id 2 = %d, want 2", got)
	}
	if got := idsToBitmask([]int{1, 3}); got != 5 {
		t.Fatalf("ids 1,3 = %d, want 5", got)
	}
	if got := idsToBitmask([]int{0, 33, -1}); got != 0 {
		t.Fatalf("out of range = %d, want 0", got)
	}
}

func TestBindCommand(t *testing.T) {
	cfg := config{Keys: []string{"F5", "F6", "F7"}}
	got := bindCommand(cfg, 7, 24)
	want := `bind "F5" "tv_listen_voice_indices 7; tv_listen_voice_indices_h 7"; bind "F6" "tv_listen_voice_indices 24; tv_listen_voice_indices_h 24"; bind "F7" "tv_listen_voice_indices -1; tv_listen_voice_indices_h -1"`
	if got != want {
		t.Fatalf("bindCommand() =\n%s\nwant\n%s", got, want)
	}
}

func TestPlaydemoName(t *testing.T) {
	if got := playdemoName(`/tmp/faceit-123.dem`); got != "faceit-123" {
		t.Fatalf("got %q", got)
	}
}

func TestPlaydemoCommand(t *testing.T) {
	folder := `/games/csgo`
	if got := playdemoCommand(folder, `/games/csgo/faceit.dem`); got != "playdemo faceit" {
		t.Fatalf("root=%q", got)
	}
	if got := playdemoCommand(folder, `/games/csgo/replays/last.dem`); got != "playdemo replays/last" {
		t.Fatalf("replay=%q", got)
	}
	if got := playdemoCommand(folder, `/tmp/upload.dem`); got != "playdemo upload" {
		t.Fatalf("outside=%q", got)
	}
}

func TestConsoleLine(t *testing.T) {
	got := consoleLine(`bind "F5" "x"`, "playdemo foo")
	if got != `bind "F5" "x"; playdemo foo` {
		t.Fatalf("got %q", got)
	}
	if got := consoleLine(`bind "F5" "x"`, ""); got != `bind "F5" "x"` {
		t.Fatalf("bind only=%q", got)
	}
	if got := consoleLine("", "playdemo foo"); got != "playdemo foo" {
		t.Fatalf("play only=%q", got)
	}
}

func TestCS2ExecutableFromCsgoFolder(t *testing.T) {
	folder := filepath.Join("Steam", "steamapps", "common", "Counter-Strike Global Offensive", "game", "csgo")
	got := cs2ExecutableFor("windows", folder)
	want := filepath.Join("Steam", "steamapps", "common", "Counter-Strike Global Offensive", "game", "bin", "win64", "cs2.exe")
	if got != want {
		t.Fatalf("windows=%q want %q", got, want)
	}
	got = cs2ExecutableFor("darwin", folder)
	want = filepath.Join("Steam", "steamapps", "common", "Counter-Strike Global Offensive", "game", "bin", "osx64", "cs2")
	if got != want {
		t.Fatalf("darwin=%q want %q", got, want)
	}
}

func TestParseListIndex(t *testing.T) {
	if parseListIndex("2", 3) != 1 {
		t.Fatal("expected 1")
	}
	if parseListIndex("0", 3) != -1 || parseListIndex("abc", 3) != -1 || parseListIndex("1", 0) != -1 {
		t.Fatal("expected -1")
	}
}

func TestIsDemoFilename(t *testing.T) {
	if !isDemoFilename("a.dem") || !isDemoFilename("a.DEM.ZST") {
		t.Fatal("expected demo filenames to be accepted")
	}
	if isDemoFilename("notes.txt") || isDemoFilename("a.zst") {
		t.Fatal("expected non-demo filenames to be rejected")
	}
}

func TestCopyDemoToGameFolderSkipsWhenAlreadyThere(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "match.dem")
	if err := os.WriteFile(src, []byte("demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config{GameFolder: dir}
	got, err := copyDemoToGameFolder(cfg, src)
	if err != nil {
		t.Fatal(err)
	}
	if got != src {
		t.Fatalf("got %q want %q", got, src)
	}
}

func TestCopyDemoToGameFolderMoveLeavesFileAlreadyInFolder(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "match.dem")
	if err := os.WriteFile(src, []byte("demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := copyDemoToGameFolder(config{GameFolder: dir, MoveDemo: true}, src)
	if err != nil {
		t.Fatal(err)
	}
	if got != src {
		t.Fatalf("got %q want %q", got, src)
	}
	if _, err := os.Stat(src); err != nil {
		t.Fatal("move must not delete a demo that is already in the game folder")
	}
}

func TestCopyDemoToGameFolderSkipsExistingDest(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "match.dem")
	if err := os.WriteFile(dst, []byte("already"), 0o644); err != nil {
		t.Fatal(err)
	}
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "match.dem")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := copyDemoToGameFolder(config{GameFolder: dir}, src)
	if err != nil {
		t.Fatal(err)
	}
	if got != dst {
		t.Fatalf("got %q want %q", got, dst)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "already" {
		t.Fatalf("overwrote dest: %q", data)
	}
}

func TestCopyDemoToGameFolderCopiesWhenMissing(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(t.TempDir(), "match.dem")
	if err := os.WriteFile(src, []byte("demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := copyDemoToGameFolder(config{GameFolder: dir}, src)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "match.dem")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "demo" {
		t.Fatalf("got %q", data)
	}
	if _, err := os.Stat(src); err != nil {
		t.Fatal("copy should leave the source file")
	}
}

func TestCopyDemoToGameFolderMovesWhenMissing(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(t.TempDir(), "match.dem")
	if err := os.WriteFile(src, []byte("demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := copyDemoToGameFolder(config{GameFolder: dir, MoveDemo: true}, src)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "match.dem")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "demo" {
		t.Fatalf("got %q", data)
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("move should remove the source file")
	}
}

func TestCopyDemoToGameFolderMoveRemovesSourceIfDestExists(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "match.dem")
	if err := os.WriteFile(dst, []byte("already"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(t.TempDir(), "match.dem")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := copyDemoToGameFolder(config{GameFolder: dir, MoveDemo: true}, src)
	if err != nil {
		t.Fatal(err)
	}
	if got != dst {
		t.Fatalf("got %q want %q", got, dst)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "already" {
		t.Fatalf("overwrote dest: %q", data)
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("move should drop the leftover source when dest already exists")
	}
}

func TestParseConfigMoveDemo(t *testing.T) {
	if parseConfig([]byte(`{"gameFolder":"C:\\cs2"}`)).MoveDemo {
		t.Fatal("expected MoveDemo to default to false")
	}
	if !parseConfig([]byte(`{"moveDemo":true}`)).MoveDemo {
		t.Fatal("expected MoveDemo true to be kept")
	}
}
