package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var libraryPathRe = regexp.MustCompile(`(?i)"path"\s+"([^"]+)"`)

var detectCS2FoldersFn = detectCS2Folders

func detectCS2Folders() []string {
	seen := make(map[string]struct{})
	found := []string{}

	add := func(p string) {
		if p == "" {
			return
		}
		clean := filepath.Clean(p)
		info, err := os.Stat(clean)
		if err != nil || !info.IsDir() {
			return
		}
		if _, ok := seen[clean]; ok {
			return
		}
		seen[clean] = struct{}{}
		found = append(found, clean)
	}

	for _, lib := range steamLibraries() {
		add(filepath.Join(lib, "steamapps", "common", "Counter-Strike Global Offensive", "game", "csgo"))
	}
	for _, p := range extraCS2Guesses() {
		add(p)
	}
	return found
}

func steamLibraries() []string {
	seen := make(map[string]struct{})
	var libs []string

	add := func(p string) {
		if p == "" {
			return
		}
		clean := filepath.Clean(unescapeVDFPath(p))
		info, err := os.Stat(clean)
		if err != nil || !info.IsDir() {
			return
		}
		if _, ok := seen[clean]; ok {
			return
		}
		seen[clean] = struct{}{}
		libs = append(libs, clean)
	}

	for _, root := range steamRoots() {
		add(root)
		for _, rel := range []string{
			filepath.Join("steamapps", "libraryfolders.vdf"),
			filepath.Join("config", "libraryfolders.vdf"),
		} {
			data, err := os.ReadFile(filepath.Join(root, rel))
			if err != nil {
				continue
			}
			for _, p := range parseLibraryFolders(string(data)) {
				add(p)
			}
		}
	}
	return libs
}

func steamRoots() []string {
	var roots []string
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "windows":
		if pf86 := os.Getenv("PROGRAMFILES(X86)"); pf86 != "" {
			roots = append(roots, filepath.Join(pf86, "Steam"))
		}
		if pf := os.Getenv("PROGRAMFILES"); pf != "" {
			roots = append(roots, filepath.Join(pf, "Steam"))
		}
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			roots = append(roots, filepath.Join(local, "Steam"))
		}
		if home != "" {
			roots = append(roots, filepath.Join(home, "Steam"))
		}
	case "darwin":
		if home != "" {
			roots = append(roots, filepath.Join(home, "Library", "Application Support", "Steam"))
		}
	default:
		if home != "" {
			roots = append(roots,
				filepath.Join(home, ".steam", "steam"),
				filepath.Join(home, ".steam", "root"),
				filepath.Join(home, ".local", "share", "Steam"),
				filepath.Join(home, ".var", "app", "com.valvesoftware.Steam", ".local", "share", "Steam"),
			)
		}
	}

	var existing []string
	seen := make(map[string]struct{})
	for _, r := range roots {
		clean := filepath.Clean(r)
		info, err := os.Stat(clean)
		if err != nil || !info.IsDir() {
			continue
		}
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		existing = append(existing, clean)
	}
	return existing
}

func extraCS2Guesses() []string {
	var guesses []string
	home, _ := os.UserHomeDir()
	csgo := filepath.Join("steamapps", "common", "Counter-Strike Global Offensive", "game", "csgo")

	if home != "" {
		guesses = append(guesses, filepath.Join(home, "Steam", csgo))
	}

	if runtime.GOOS != "windows" {
		return guesses
	}

	for _, letter := range []string{"C", "D", "E", "F", "G", "H"} {
		guesses = append(guesses,
			letter+`:\`+csgo,
			letter+`:\Steam\`+csgo,
			letter+`:\SteamLibrary\`+csgo,
			letter+`:\Games\Steam\`+csgo,
			letter+`:\Games\SteamLibrary\`+csgo,
			letter+`:\Program Files (x86)\Steam\`+csgo,
			letter+`:\Program Files\Steam\`+csgo,
		)
	}
	return guesses
}

func parseLibraryFolders(content string) []string {
	matches := libraryPathRe.FindAllStringSubmatch(content, -1)
	var paths []string
	seen := make(map[string]struct{})
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		p := unescapeVDFPath(m[1])
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		paths = append(paths, p)
	}
	return paths
}

func unescapeVDFPath(p string) string {
	p = strings.ReplaceAll(p, `\\`, `\`)
	return strings.TrimSpace(p)
}

func looksLikeCSGOFolder(path string) bool {
	norm := strings.ReplaceAll(filepath.Clean(path), `\`, "/")
	norm = strings.TrimRight(strings.ToLower(norm), "/")
	if strings.HasSuffix(norm, "/game/csgo") {
		return true
	}
	if _, err := os.Stat(filepath.Join(path, "cfg")); err == nil {
		return true
	}
	return false
}

func normalizeGameFolder(path string) string {
	path = filepath.Clean(strings.TrimSpace(path))
	if path == "" || path == "." {
		return path
	}
	if looksLikeCSGOFolder(path) {
		return path
	}
	for _, rel := range []string{
		filepath.Join("game", "csgo"),
		"csgo",
		filepath.Join("Counter-Strike Global Offensive", "game", "csgo"),
	} {
		cand := filepath.Join(path, rel)
		info, err := os.Stat(cand)
		if err == nil && info.IsDir() {
			return cand
		}
	}
	return path
}
