package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

// Ask question to user and return inout value
func Ask(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("> %v: ", question)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(input)
	return value, nil
}
