package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	GameFolder   string   `json:"gameFolder"`
	Keys         []string `json:"keys"`
	CheckUpdates bool     `json:"checkUpdates"`
	MoveDemo     bool     `json:"moveDemo"`
}

func defaultConfig() config {
	return config{
		Keys:         []string{"F5", "F6", "F7"},
		CheckUpdates: true,
	}
}

func configPath() (string, error) {
	return configPathFn()
}

var configPathFn = defaultConfigPath

func defaultConfigPath() (string, error) {
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

	return parseConfig(data)
}

func parseConfig(data []byte) config {
	cfg := defaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig()
	}
	if len(cfg.Keys) != 3 {
		cfg.Keys = defaultConfig().Keys
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
