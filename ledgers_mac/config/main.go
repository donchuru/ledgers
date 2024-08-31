package main

import (
	"fmt"
	"os"
	"bufio"
	"path/filepath"
)

// error checker
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main(){
	var name string
	var location string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Prompt 1/2: Please enter your name: ")
	if scanner.Scan() {
		name = scanner.Text()
		// fmt.Printf("Input was: %q\n", name)
	}

	fmt.Println("Prompt 2/2: Please enter where you'd like to store your journals (Paste with Command+V): ")
	if scanner.Scan() {
		location = scanner.Text()
		// fmt.Printf("Input was: %q\n", location)
	}

	// Resolve the home directory
	homeDir, err := os.UserHomeDir()
	check(err)

	// Define the configuration file path
	configDir := filepath.Join(homeDir, ".ledgers_config")
	configFilePath := filepath.Join(configDir, "init.txt")

	// Ensure the directory exists
	err = os.MkdirAll(configDir, 0755)
	check(err)

	d1 := []byte(name + "\n" + location)
    err2 := os.WriteFile(configFilePath, d1, 0644)
    check(err2)
}