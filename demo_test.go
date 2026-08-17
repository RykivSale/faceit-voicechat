package main

import "testing"

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
