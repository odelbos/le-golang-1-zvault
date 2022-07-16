package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	// Get default config absolute path and file name.
	defaultConfig, err := DefaultConfig()
	if err != nil {
		fmt.Println("Cannot get default config pathfile !")
		os.Exit(1)
	}

	// Define application flags and commands
	app := cli.NewApp()
	app.Name = "zvault"
	app.Usage = "store/restore files in encrypted vault"
	app.HelpName = "zvault"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: defaultConfig,
			Usage: "The configuration file to use.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize a new vault",
			Action:  InitCmd,
		},
		{
			Name:      "put",
			Aliases:   []string{"p"},
			Usage:     "Store a file in vault",
			Action:    PutCmd,
			ArgsUsage: "</path/to/file>",
		},
		{
			Name:      "get",
			Aliases:   []string{"g"},
			Usage:     "Restore a file from vault",
			ArgsUsage: "<id>",
			Action:    GetCmd,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
