package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseLibraryFolders(t *testing.T) {
	vdf := `"libraryfolders"
{
	"0"
	{
		"path"		"C:\\Program Files (x86)\\Steam"
	}
	"1"
	{
		"path"		"D:\\SteamLibrary"
	}
}`
	got := parseLibraryFolders(vdf)
	if len(got) != 2 {
		t.Fatalf("len=%d want 2: %#v", len(got), got)
	}
	if got[0] != `C:\Program Files (x86)\Steam` {
		t.Fatalf("got[0]=%q", got[0])
	}
	if got[1] != `D:\SteamLibrary` {
		t.Fatalf("got[1]=%q", got[1])
	}
}

func TestLooksLikeCSGOFolder(t *testing.T) {
	p := `C:\Program Files (x86)\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo`
	if !looksLikeCSGOFolder(p) {
		t.Fatalf("expected %q to look like csgo folder", p)
	}
}

func TestNormalizeGameFolderWalksIntoCSGO(t *testing.T) {
	root := t.TempDir()
	csgo := filepath.Join(root, "Counter-Strike Global Offensive", "game", "csgo")
	if err := os.MkdirAll(csgo, 0o755); err != nil {
		t.Fatal(err)
	}
	got := normalizeGameFolder(root)
	if got != csgo {
		t.Fatalf("got %q want %q", got, csgo)
	}
}
