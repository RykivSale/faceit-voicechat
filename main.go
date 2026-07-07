package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

type config struct {
	GameFolder string   `json:"gameFolder"`
	Keys       []string `json:"keys"`
}

func defaultConfig() config {
	return config{Keys: []string{"F5", "F6", "F7"}}
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "faceit-voicechat", "config.json"), nil
}

func loadConfig() config {
	cfg := defaultConfig()

	path, err := configPath()
	if err != nil {
		return cfg
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	var loaded config
	if err := json.Unmarshal(data, &loaded); err != nil {
		return cfg
	}

	if loaded.GameFolder != "" {
		cfg.GameFolder = loaded.GameFolder
	}
	if len(loaded.Keys) == 3 {
		cfg.Keys = loaded.Keys
	}

	return cfg
}

func saveConfig(cfg config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var rawPath string
	if len(os.Args) >= 2 {
		rawPath = os.Args[1]
	} else {
		fmt.Print("Enter path to .dem or .dem.zst file: ")
		line, _ := reader.ReadString('\n')
		rawPath = strings.Trim(strings.TrimSpace(line), `"`)

		if rawPath == "" {
			fmt.Println("No file path given.")
			waitForExit()
			return
		}
	}

	demoPath, cleanup, err := resolveDemoPath(rawPath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		waitForExit()
		return
	}
	if cleanup != nil {
		defer cleanup()
	}

	ctMask, tMask, err := parseDemo(demoPath)
	if err != nil {
		fmt.Println("Failed to parse demo:", err)
		waitForExit()
		return
	}

	cfg := loadConfig()

	for {
		fmt.Println()
		fmt.Println("Press Enter - get bind")
		fmt.Println("Press S - settings")
		fmt.Println("Press Q - quit")
		fmt.Print("> ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch strings.ToLower(line) {
		case "":
			printBind(cfg, ctMask, tMask)
			copyDemoAndPrintCommand(cfg, demoPath)
		case "s":
			settingsMenu(reader, &cfg)
		case "q":
			return
		default:
			fmt.Println("Unknown option.")
		}
	}
}

// resolveDemoPath returns a path to a plain .dem file, transparently decompressing
// rawPath first if it's a .zst archive. The returned cleanup func (if non-nil) removes
// the temporary decompressed file and must be called once the caller is done with it.
func resolveDemoPath(rawPath string) (path string, cleanup func(), err error) {
	if !strings.EqualFold(filepath.Ext(rawPath), ".zst") {
		return rawPath, nil, nil
	}

	fmt.Println("Decompressing .zst archive...")

	in, err := os.Open(rawPath)
	if err != nil {
		return "", nil, err
	}
	defer in.Close()

	dec, err := zstd.NewReader(in)
	if err != nil {
		return "", nil, err
	}
	defer dec.Close()

	tmpDir, err := os.MkdirTemp("", "faceit-voicechat-*")
	if err != nil {
		return "", nil, err
	}

	outPath := filepath.Join(tmpDir, strings.TrimSuffix(filepath.Base(rawPath), filepath.Ext(rawPath)))
	out, err := os.Create(outPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, dec); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}

	return outPath, func() { os.RemoveAll(tmpDir) }, nil
}

func parseDemo(demoPath string) (ctMask uint32, tMask uint32, err error) {
	f, err := os.Open(demoPath)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()

	var ctIDs []int
	var tIDs []int

	p.RegisterEventHandler(func(e events.MatchStart) {
		participants := p.GameState().Participants()
		cts := participants.TeamMembers(common.TeamCounterTerrorists)
		ts := participants.TeamMembers(common.TeamTerrorists)

		for _, ct := range cts {
			ctIDs = append(ctIDs, ct.EntityID)
		}
		for _, t := range ts {
			tIDs = append(tIDs, t.EntityID)
		}
	})

	if err := p.ParseToEnd(); err != nil {
		return 0, 0, err
	}

	return idsToBitmask(ctIDs), idsToBitmask(tIDs), nil
}

func idsToBitmask(ids []int) uint32 {
	var mask uint32
	for _, id := range ids {
		if id >= 1 && id <= 32 {
			mask |= 1 << (id - 1)
		}
	}
	return mask
}

func printBind(cfg config, ctMask, tMask uint32) {
	fmt.Println()
	fmt.Printf(
		"bind \"%s\" \"tv_listen_voice_indices %d; tv_listen_voice_indices_h %d\"; bind \"%s\" \"tv_listen_voice_indices %d; tv_listen_voice_indices_h %d\"; bind \"%s\" \"tv_listen_voice_indices -1; tv_listen_voice_indices_h -1\"\n",
		cfg.Keys[0], ctMask, ctMask,
		cfg.Keys[1], tMask, tMask,
		cfg.Keys[2],
	)
}

func copyDemoAndPrintCommand(cfg config, demoPath string) {
	base := filepath.Base(demoPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	if cfg.GameFolder == "" {
		fmt.Println()
		fmt.Println("Game folder isn't set, so the demo wasn't copied (press S -> 1 to set it).")
		fmt.Printf("playdemo %s\n", name)
		return
	}

	dst := filepath.Join(cfg.GameFolder, base)
	if err := copyFile(demoPath, dst); err != nil {
		fmt.Println()
		fmt.Println("Failed to copy demo to game folder:", err)
		return
	}

	fmt.Println()
	fmt.Printf("Demo copied to: %s\n", dst)
	fmt.Printf("playdemo %s\n", name)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func settingsMenu(reader *bufio.Reader, cfg *config) {
	for {
		fmt.Println()
		fmt.Println("--- Settings ---")
		fmt.Printf("1) Set game folder (current: %s)\n", displayOrNone(cfg.GameFolder))
		fmt.Printf("2) Change keybinds (current: %s, %s, %s)\n", cfg.Keys[0], cfg.Keys[1], cfg.Keys[2])
		fmt.Println("3) Back")
		fmt.Print("> ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter full path to your game folder (e.g. ...\\Counter-Strike Global Offensive\\game\\csgo): ")
			path, _ := reader.ReadString('\n')
			path = strings.Trim(strings.TrimSpace(path), `"`)

			info, err := os.Stat(path)
			if err != nil || !info.IsDir() {
				fmt.Println("That folder doesn't exist. Not saved.")
				continue
			}

			cfg.GameFolder = path
			if err := saveConfig(*cfg); err != nil {
				fmt.Println("Failed to save settings:", err)
				continue
			}
			fmt.Println("Saved.")
		case "2":
			fmt.Print("Enter 3 keys separated by commas (e.g. F5,F6,F7): ")
			line, _ := reader.ReadString('\n')
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				fmt.Println("Please enter exactly 3 keys separated by commas.")
				continue
			}

			keys := make([]string, 0, 3)
			for _, p := range parts {
				key := strings.ToUpper(strings.TrimSpace(p))
				if key == "" {
					keys = nil
					break
				}
				keys = append(keys, key)
			}
			if len(keys) != 3 {
				fmt.Println("Please enter exactly 3 non-empty keys separated by commas.")
				continue
			}

			cfg.Keys = keys
			if err := saveConfig(*cfg); err != nil {
				fmt.Println("Failed to save settings:", err)
				continue
			}
			fmt.Println("Saved.")
		case "3", "b", "":
			return
		default:
			fmt.Println("Unknown option.")
		}
	}
}

func displayOrNone(s string) string {
	if s == "" {
		return "not set"
	}
	return s
}

func waitForExit() {
	fmt.Println()
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
