package main

import (
	"os"
	"path/filepath"
)

const (
	CONFIG_DIR = ".config"
	CONFIG_NAME = "zvault.json"
)

func DefaultConfig() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userHome, CONFIG_DIR, CONFIG_NAME), nil
}
