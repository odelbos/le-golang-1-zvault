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

func ConfigExists(config string) (bool, error) {
	if _, err := os.Stat(config); err != nil {
		return false, err
	}
	return true, nil
}
