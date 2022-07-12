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
		// The configuration file already exists
		//
		// TODO : Do you want to overwrite ?
		// 
		fmt.Println("The config file already exists !!!")
	}
}
