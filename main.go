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

	fmt.Printf("Conf file : %v\n", configPath)
	fmt.Printf("Args : %v\n", args)

	if len(args) < 1 {
		fmt.Println("You must provide a command: init, put, get")
		os.Exit(1)
	}
	cmd := args[0]

	// ---------------------------------------------------
	// Init Command
	// ---------------------------------------------------
	if cmd == "init" {
		fmt.Println("We got an 'init' command")
		os.Exit(0)
	}

	// ---------------------------------------------------
	// Put Command
	// ---------------------------------------------------
	if cmd == "put" {
		fmt.Println("We got a 'put' command")
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
