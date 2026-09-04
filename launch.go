package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const cs2SteamAppID = "730"

var launchCS2Fn = launchCS2

func launchCS2(gameFolder string) error {
	steamErr := openSteamApp(cs2SteamAppID)
	if steamErr == nil {
		return nil
	}
	exe := cs2Executable(gameFolder)
	if exe == "" {
		return steamErr
	}
	if _, err := os.Stat(exe); err != nil {
		return steamErr
	}
	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	return cmd.Start()
}

func openSteamApp(appID string) error {
	u := "steam://rungameid/" + appID
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/C", "start", "", u).Start()
	case "darwin":
		return exec.Command("open", u).Start()
	default:
		if err := exec.Command("xdg-open", u).Start(); err == nil {
			return nil
		}
		return exec.Command("steam", "-applaunch", appID).Start()
	}
}

func cs2Executable(gameFolder string) string {
	return cs2ExecutableFor(runtime.GOOS, gameFolder)
}

func cs2ExecutableFor(goos, gameFolder string) string {
	if strings.TrimSpace(gameFolder) == "" {
		return ""
	}
	gameDir := filepath.Clean(gameFolder)
	if strings.EqualFold(filepath.Base(gameDir), "csgo") {
		gameDir = filepath.Dir(gameDir)
	}
	switch goos {
	case "windows":
		return filepath.Join(gameDir, "bin", "win64", "cs2.exe")
	case "darwin":
		return filepath.Join(gameDir, "bin", "osx64", "cs2")
	default:
		return filepath.Join(gameDir, "bin", "linuxsteamrt64", "cs2")
	}
}
