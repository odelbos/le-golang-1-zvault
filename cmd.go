package main

import (
	"fmt"
	"os"
)

func InitCmd(configPath string) {
	fmt.Printf("InitCmd : Configuration file : %v\n", configPath)

	exists, err := ConfigExists(configPath)
	if err != nil {
		fmt.Println("Please run: zvault init")
		os.Exit(1)
	}

	if exists {
		fmt.Println("The config file already exists !")

		answer, err := Ask("Do you want to overwrite it? (yes/no)")
		if err != nil {
			fmt.Println("Caannot read user input !")
			os.Exit(1)
		}

		if answer != "yes" {
			fmt.Println("-quit-")
			os.Exit(0)
		}
		fmt.Println("")        // Output empty line
	}

	// Ask for storage folders and master pawwsord
	//
	dataPath, err := Ask("Data folder")
	if err != nil {
		fmt.Println("Caannot read user input !")
		os.Exit(1)
	}
	filesPath, err := Ask("Files folder")
	if err != nil {
		fmt.Println("Caannot read user input !")
		os.Exit(1)
	}
	// Ask password
	pwd, err := AskPwdTwice()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create new config and save it
	config := NewConfig(dataPath, filesPath)
	config.Save(configPath, pwd)
	fmt.Println("Vault configuration saved.")
}

func PutCmd(config *Config, args []string) {
	if len(args) != 1 {
		fmt.Println("With the 'put' command you must provide the path of the file to store in the vault !")
		os.Exit(1)
	}

	filePath := args[0]
	fileId, err := config.Vault.Put(filePath)
	if err != nil {
		fmt.Println("Cannot store file in vault !")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("File stored, id: %v\n", fileId)
}

func GetCmd(config *Config, args []string) {
	if len(args) != 1 {
		fmt.Println("With the 'get' command you must provide the file id !")
		os.Exit(1)
	}

	fileId := args[0]
	err := config.Vault.Get(fileId)
	if err != nil {
		fmt.Println("Cannot get the file !")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("File restored.")
}
