package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	VERSION     = 1
	CONFIG_DIR  = ".config"
	CONFIG_NAME = "zvault.json"
)

type Config struct {
	Version uint
	Vault   Vault
}

type EncryptedConfig struct {
	Version        uint   `json:"r"`
	Salt           []byte `json:"s"`
	Iter           int    `json:"t"`
	Iv             []byte `json:"i"`
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
	_, err := os.Stat(config)
	if err == nil {
		return true, err
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
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
	masterKey := GenCryptoRand(AES_KEY_SIZE)
	vault := Vault{
		DataPath:  absDataPath,
		FilesPath: absFilesPath,
		MasterKey: masterKey,
	}
	config := Config{
		Version: VERSION,
		Vault:   vault,
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
		Version:        c.Version,
		Salt:           salt,
		Iter:           PBKDF2_ITER,
		Iv:             *iv,
		EncryptedVault: *encryptedVault,
	}
	jsonEncryptedConfig, err := json.MarshalIndent(eConfig, "", "  ")
	if err != nil {
		return err
	}
	jsonSize := len(jsonEncryptedConfig)

	// Save config
	err = os.WriteFile(configPath, jsonEncryptedConfig[:jsonSize], 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig(configPath string, pwd []byte) (Config, error) {
	// Load the encrypted config
	jsonEncryptedConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Cannot read configuration file !")
		os.Exit(1)
	}
	var eConfig EncryptedConfig
	err = json.Unmarshal(jsonEncryptedConfig, &eConfig)
	if err != nil {
		return Config{}, err
	}

	// Get back the derived key.
	derivedKey := GenPBKDF2WithSalt(pwd, eConfig.Salt)

	// Descrupt the vault configuration
	decryptedVault, err := Decrypt(&eConfig.EncryptedVault, &derivedKey, &eConfig.Iv)
	if err != nil {
		return Config{}, err
	}
	var vault Vault
	err = json.Unmarshal(*decryptedVault, &vault)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Version: eConfig.Version,
		Vault:   vault,
	}, nil
}
