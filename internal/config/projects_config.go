package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ProjectConfig struct {
	DefaultLabel    string `json:"default_label"`
	DefaultReviewer string `json:"default_reviewer"`
}

func ensureConfig() (string, error) {
	dir, err := os.UserConfigDir()

	if err != nil {
		return "", fmt.Errorf("reading config dir: %w", err)
	}

	configDir := filepath.Join(dir, "gh-targetprocess")
	configFile := filepath.Join(configDir, "projects.json")

	if _, err := os.Stat(configDir); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		if err := os.Mkdir(configDir, 0700); err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(configFile); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		f, err := os.Create(configFile)

		if err != nil {
			return "", err
		}

		defer f.Close()
	}

	return configFile, nil
}

func Read() (*ProjectConfig, error) {
	filepath, err := ensureConfig()

	if err != nil {
		return nil, fmt.Errorf("ensuring config: %w", err)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file (%s): %w", filepath, err)
	}

	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &config, nil
}
