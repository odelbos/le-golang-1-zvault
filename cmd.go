package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func InitCmd(c *cli.Context) error {
	configPath := c.GlobalString("config")
	fmt.Printf("Configuration file : %v\n", configPath)

	exists, err := ConfigExists(configPath)
	if err != nil {
		return cli.NewExitError("Cannot verify if config file exists!", 1)
	}

	if exists {
		fmt.Println("The config file already exists !")
		answer, err := Ask("Do you want to overwrite it? (yes/no)")
		if err != nil {
			return cli.NewExitError("Cannot read user input!", 1)
		}

		if answer != "yes" {
			fmt.Println("-quit-")
			return nil
		}
	}

	// Ask for storage folders and master pawwsord
	dataPath, err := Ask("Data folder")
	if err != nil {
		return cli.NewExitError("Cannot read user input!", 1)
	}
	filesPath, err := Ask("Files folder")
	if err != nil {
		return cli.NewExitError("Cannot read user input!", 1)
	}
	pwd, err := AskPwdTwice()
	if err != nil {
		return err
	}

	// Create and save the configuration file
	config := NewConfig(dataPath, filesPath)
	config.Save(configPath, pwd)
	fmt.Println("Vault configuration created.")
	return nil
}

func PutCmd(c *cli.Context) error {
	// Check if file is provided
	if len(c.Args()) != 1 {
		return cli.NewExitError("You must provide the file to store!", 1)
	}

	// Load config
	config, err := askPwdAndLoadConfig(c)
	if err != nil {
		return cli.NewExitError("Caannot load config!", 1)
	}

	filePath := c.Args().First()
	fileId, err := config.Vault.Put(filePath)
	if err != nil {
		return cli.NewExitError("cannot store file in vault", 1)
	}
	fmt.Printf("File stored, id: %v\n", fileId)
	return nil
}

func GetCmd(c *cli.Context) error {
	// Check if file is provided
	if len(c.Args()) != 1 {
		return cli.NewExitError("You must provide the id of file to restore!", 1)
	}

	// Load config
	config, err := askPwdAndLoadConfig(c)
	if err != nil {
		return cli.NewExitError("Caannot load config!", 1)
	}

	fileId := c.Args().First()
	fileName, err := config.Vault.Get(fileId)
	if err != nil {
		return cli.NewExitError("Caannot read file!", 1)
	}
	fmt.Printf("File restored, name: %v\n", fileName)
	return nil
}

// ----------------------------------------------------
// Helper functions
// ----------------------------------------------------
func askPwdAndLoadConfig(c *cli.Context) (Config, error) {
	configPath := c.GlobalString("config")

	pwd, err := AskPwd()
	if err != nil {
		return Config{}, cli.NewExitError("Caannot read password!", 1)
	}
	config, err := LoadConfig(configPath, pwd)
	if err != nil {
		return Config{}, cli.NewExitError("Caannot load configuration", 1)
	}
	return config, nil
}
