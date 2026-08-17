package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	cli := false
	var rawPath string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--cli":
			cli = true
		case "-h", "--help":
			printUsage()
			return
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Println("Unknown flag:", arg)
				printUsage()
				os.Exit(2)
			}
			if rawPath == "" {
				rawPath = arg
			}
		}
	}

	if cli {
		runCLI(rawPath)
		return
	}
	runWeb(rawPath)
}

func printUsage() {
	fmt.Println("faceit-voicechat — generate CS2 voice listen binds from a demo")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  faceit-voicechat [demo.dem|demo.dem.zst]")
	fmt.Println("  faceit-voicechat --cli [demo.dem|demo.dem.zst]")
	fmt.Println()
	fmt.Println("By default a local page opens in your browser.")
	fmt.Println("Use --cli for the old console menu.")
}

func runCLI(rawPath string) {
	reader := bufio.NewReader(os.Stdin)

	if rawPath == "" {
		fmt.Print("Enter path to .dem or .dem.zst file: ")
		line, _ := reader.ReadString('\n')
		rawPath = strings.Trim(strings.TrimSpace(line), `"`)

		if rawPath == "" {
			fmt.Println("No file path given.")
			waitForExit()
			return
		}
	}

	if strings.HasSuffix(strings.ToLower(rawPath), ".zst") {
		fmt.Println("Decompressing .zst archive...")
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

func printBind(cfg config, ctMask, tMask uint32) {
	fmt.Println()
	fmt.Println(bindCommand(cfg, ctMask, tMask))
}

func copyDemoAndPrintCommand(cfg config, demoPath string) {
	name := playdemoName(demoPath)

	if cfg.GameFolder == "" {
		fmt.Println()
		fmt.Println("Game folder isn't set, so the demo wasn't copied (press S -> 1 to set it).")
		fmt.Printf("playdemo %s\n", name)
		return
	}

	dst, err := copyDemoToGameFolder(cfg, demoPath)
	if err != nil {
		fmt.Println()
		fmt.Println("Failed to copy demo to game folder:", err)
		return
	}

	fmt.Println()
	fmt.Printf("Demo copied to: %s\n", dst)
	fmt.Printf("playdemo %s\n", name)
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
			found := detectCS2Folders()
			if len(found) > 0 {
				fmt.Println("Found CS2 folder(s):")
				for i, p := range found {
					fmt.Printf("  %d) %s\n", i+1, p)
				}
				fmt.Print("Enter a number, or paste a full path (e.g. ...\\Counter-Strike Global Offensive\\game\\csgo): ")
			} else {
				fmt.Print("Enter full path to your game folder (e.g. ...\\Counter-Strike Global Offensive\\game\\csgo): ")
			}
			path, _ := reader.ReadString('\n')
			path = strings.Trim(strings.TrimSpace(path), `"`)
			if n := parseListIndex(path, len(found)); n >= 0 {
				path = found[n]
			}

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

func parseListIndex(s string, n int) int {
	if n == 0 || s == "" {
		return -1
	}
	var idx int
	if _, err := fmt.Sscanf(s, "%d", &idx); err != nil {
		return -1
	}
	if idx < 1 || idx > n {
		return -1
	}
	return idx - 1
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
