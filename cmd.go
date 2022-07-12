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
	}

	//
	// TODO : Ask for storage folders and master pawwsord
	//
}
