package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	// "sort"
)

var location string

// error checker
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// check if file exists
func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// append to existing file
func appendToFile(filepath string, content string) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		check(err)
		return
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		check(err)
	}
}

type LedgerEntry struct {
	filename string
	tags     []string
	content  string
}

func createLedgerEntry(args []string) LedgerEntry {
	entry := LedgerEntry{
		content: time.Now().Format("Monday, Jan 2, 2006") + "\n\n",
	}

	// Check for -t flag and extract tags
	tIndex := findIndex(args, "-t")
	if tIndex != -1 && tIndex < len(args)-1 {
		entry.tags = args[tIndex+1:]
		args = args[:tIndex] // Remove -t and tags from args
	}

	// Set filename based on arguments
	if len(args) <= 1 {
		entry.filename = time.Now().Format("2006-01-02") + ".txt"
	} else {
		entry.filename = strings.Join(args[1:], " ") + ".txt"
	}

	return entry
}

func (entry LedgerEntry) save(location string) error {
	filepath := filepath.Join(location, entry.filename)
	
	// Prepare content with tags if present
	content := entry.content
	if len(entry.tags) > 0 {
		content = "tags: " + strings.Join(entry.tags, ", ") + "\n" + entry.filename + "\n" + content
	}

	// Append or create file
	if fileExists(filepath) {
		separator := "\n\n"
		if len(entry.tags) > 0 {
			separator += "tags: " + strings.Join(entry.tags, ", ") + "\n"
		}
		appendToFile(filepath, separator + entry.content)
	} else {
		err := os.WriteFile(filepath, []byte(content), 0644)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func main() {
	/* take in command line arguments
	User inputs:
		ledger -> make a new ledger named today's date
		ledger "new Doc"  -> make a new ledger named new Doc
		ledger "new Doc" -t life -> make a new ledger named new Doc with the life tag
	*/

	// fetch location where all journals are stored
	homeDir, err := os.UserHomeDir()
	check(err)

	// Define the configuration file path
	configFilePath := filepath.Join(homeDir, ".ledgers_config", "init.txt")

	f, _ := os.Open(configFilePath)
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	scanner.Scan()
	location = scanner.Text()

	entry := createLedgerEntry(os.Args)
	err = entry.save(location)
	if err != nil {
		fmt.Printf("Error saving ledger: %v\n", err)
		return
	}

	fmt.Println("Don't forget to save (Ctrl + S)")
	cmd := exec.Command("open", "-a", "TextEdit", filepath.Join(location, entry.filename))
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}

// helper functions
func findIndex(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}
