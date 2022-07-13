package main

import (
	"fmt"
	"os"
	"flag"
)

func main() {
	// Get default config absolute path and file name.
	defaultConfig, err := DefaultConfig()
	if err != nil {
		fmt.Println("Cannot get default config pathfile !")
		os.Exit(1)
	}

	// Parse flags
	var configPath string
	flag.StringVar(&configPath, "c", defaultConfig, "Vault configuration file")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("You must provide a command: init, put, get")
		os.Exit(1)
	}
	cmd := args[0]

	// ---------------------------------------------------
	// Init Command
	// ---------------------------------------------------
	if cmd == "init" {
		InitCmd(configPath)
		os.Exit(0)
	}

	// All the other commands need a vault configuration.
	// Ask password and load config.
	pwd, err := AskPwd()
	if err != nil {
		fmt.Println("Caannot read password !", err)
		os.Exit(1)
	}
	config := LoadConfig(configPath, pwd)

	// ---------------------------------------------------
	// Put Command
	// ---------------------------------------------------
	if cmd == "put" {
		PutCmd(&config, args[1:])
		os.Exit(0)
	}

	// ---------------------------------------------------
	// Get Command
	// ---------------------------------------------------
	if cmd == "get" {
		fmt.Println("We got a 'get' command")
		os.Exit(0)
	}

	// ---------------------------------------------------
	// Error unknown command
	// ---------------------------------------------------
	fmt.Println("Unknown command !")
	os.Exit(1)
}
