package main

import (
	"os"
	"path/filepath"
	"encoding/json"
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

type EncryptedConfig struct {
	Version uint `json:"r"`
	Salt []byte `json:"s"`
	Iter int `json:"t"`
	Iv []byte `json:"i"`
	EncryptedVault []byte `json:"v"`
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

func NewConfig(dataPath string, filesPath string) *Config {
	// Make sure we have absolute pathes
	absDataPath, err := filepath.Abs(dataPath)
	if err != nil {
		panic("Cannot get absolute path !")
	}
	absFilesPath, err := filepath.Abs(filesPath)
	if err != nil {
		panic("Cannot get absolute path !")
	}
	// Generate a master key
	masterKey := GenCryptoRand(32)
	vault := Vault{
		DataPath: absDataPath,
		FilesPath: absFilesPath,
		MasterKey: masterKey,
	}
	config := Config{
		Version: VERSION,
		Vault:  vault,
	}
	return &config
}

func (c *Config) Save(configPath string, pwd []byte) error {
	// Encode Vault config in JSON and encrypt it with
	// a derived password key.
	jsonVault, err := json.Marshal(c.Vault)
	if err != nil {
		return err
	}
	derivedKey, salt := GenPBKDF2(pwd)
	encryptedVault, iv, err := EncryptWithKey(&jsonVault, &derivedKey)
	if err != nil {
		return err
	}
	// Create the encrypted configuration
	eConfig := EncryptedConfig{
		Version: c.Version,
		Salt: salt,
		Iter: PBKDF2_ITER,
		Iv: *iv,
		EncryptedVault: *encryptedVault,
	}
	jsonEncryptedConfig, err := json.MarshalIndent(eConfig, "", "  ")
	if err != nil {
		return err
	}
	// Save config
	err = os.WriteFile(configPath, jsonEncryptedConfig[:len(jsonEncryptedConfig)], 0644)
	if err != nil {
		return err
	}
	return nil
}
