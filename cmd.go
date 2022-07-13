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
	dataPath, err := Ask("Data path")
	if err != nil {
		fmt.Println("Caannot read user input !")
		os.Exit(1)
	}
	filesPath, err := Ask("Files path")
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

	fmt.Printf("Data path: %v\n", dataPath)
	fmt.Printf("Files path: %v\n", filesPath)
	fmt.Printf("Password: %v\n", pwd)

	//
	// TODO : save configuration
	//
	config := NewConfig(dataPath, filesPath, pwd)

	fmt.Println(config)
}
