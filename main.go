package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("You must provide a command: init, put, get")
		os.Exit(1)
	}

	fmt.Printf("You choose cmd : %v\n", args[0])
}
