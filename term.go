package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"bytes"
	"golang.org/x/crypto/ssh/terminal"
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

// Ask user for a password twice and verify equality
func AskPwdTwice() ([]byte, error) {
	fmt.Print("> Enter Password: ")
	pwd1, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, err
	}

	fmt.Print("\n> Repeat Password: ")
	pwd2, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	fmt.Println("")

	if ! bytes.Equal(pwd1, pwd2) {
		return nil, fmt.Errorf("Passwords do not match !")
	}
	return pwd1, nil
}

// Prompt for a password
func AskPwd() ([]byte, error) {
	fmt.Print("> Enter Password: ")
	pwd, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	fmt.Println("")
	return pwd, nil
}
