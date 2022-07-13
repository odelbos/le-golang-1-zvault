package main

import (
	"os"
	"path/filepath"
)

const (
	VERSION = 1
	CONFIG_DIR = ".config"
	CONFIG_NAME = "zvault.json"
)

type Config struct {
	Version uint
	Vault Vault
}

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

func NewConfig(dataPath string, filesPath string, pwd []byte) *Config {
	// Generate a master key
	masterKey := GenCryptoRand(32)
	vault := Vault{
		DataPath: dataPath,
		FilesPath: filesPath,
		MasterKey: masterKey,
	}

	config := Config{
		Version: VERSION,
		Vault:  vault,
	}
	return &config
}
