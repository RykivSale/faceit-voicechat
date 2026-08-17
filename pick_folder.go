package main

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

var errPickCancelled = errors.New("folder pick cancelled")

var pickFolderFn = nativePickFolder

func nativePickFolder() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return pickFolderWindows()
	case "darwin":
		return pickFolderDarwin()
	default:
		return pickFolderLinux()
	}
}

func pickFolderWindows() (string, error) {
	script := `
Add-Type -AssemblyName System.Windows.Forms | Out-Null
$d = New-Object System.Windows.Forms.FolderBrowserDialog
$d.Description = 'Select CS2 game\csgo folder'
$d.ShowNewFolderButton = $false
if ($d.ShowDialog() -ne [System.Windows.Forms.DialogResult]::OK) { exit 1 }
[Console]::Out.Write($d.SelectedPath)
`
	cmd := exec.Command("powershell", "-NoProfile", "-STA", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		return "", errPickCancelled
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", errPickCancelled
	}
	return path, nil
}

func pickFolderDarwin() (string, error) {
	script := `try
	POSIX path of (choose folder with prompt "CS2 folder (game/csgo)")
on error number -128
	return ""
end try`
	cmd := exec.Command("osascript", "-e", script)
	out, err := cmd.Output()
	if err != nil {
		return "", errPickCancelled
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", errPickCancelled
	}
	return path, nil
}

func pickFolderLinux() (string, error) {
	cmd := exec.Command("zenity", "--file-selection", "--directory", "--title=CS2 game/csgo")
	out, err := cmd.Output()
	if err != nil {
		return "", errPickCancelled
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", errPickCancelled
	}
	return path, nil
}
