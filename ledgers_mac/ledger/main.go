package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
	// "sort"

	"github.com/donchuru/ledgers/ledgers_mac/header"
	"github.com/donchuru/ledgers/ledgers_mac/web"
)

var location string

// error checker
func check(e error) {
	if e != nil {
		panic(e)
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

// removeFlag returns args without flag, and whether flag was there
func removeFlag(args []string, flag string) ([]string, bool) {
	i := findIndex(args, flag)
	if i == -1 {
		return args, false
	}
	return slices.Delete(slices.Clone(args), i, i+1), true
}

type LedgerEntry struct {
	filename string
	tags     []string
}

func createLedgerEntry(args []string) LedgerEntry {
	entry := LedgerEntry{}

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
	title := strings.TrimSuffix(entry.filename, ".txt")

	// New journals get a fresh header; existing ones keep theirs (or get one
	// built from the old format) and have the new entry appended
	h := header.Header{Title: title, Created: time.Now().Format(header.DateLayout)}
	body := "\n"
	info, err := os.Stat(filepath)
	if err == nil {
		raw, err := os.ReadFile(filepath)
		if err != nil {
			return err
		}
		var ok bool
		if h, body, ok = header.Parse(string(raw)); !ok {
			h, body = header.FromLegacy(title, string(raw), info.ModTime())
		}
		body += "\n\n"
	} else if !os.IsNotExist(err) {
		return err
	}
	h.AddTags(entry.tags...)

	content := h.String() + body + time.Now().Format(header.StampLayout) + "\n\n"
	return os.WriteFile(filepath, []byte(content), 0644)
}

func main() {
	/* take in command line arguments
	User inputs:
		ledger -> make a new ledger named today's date
		ledger "new Doc"  -> make a new ledger named new Doc
		ledger "new Doc" -t life -> make a new ledger named new Doc with the life tag
		ledger "new Doc" -web -> same, but write in Chrome instead of TextEdit
	*/
	web.ServeIfRequested()

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

	args, inBrowser := removeFlag(os.Args, "-web")
	entry := createLedgerEntry(args)
	err = entry.save(location)
	if err != nil {
		fmt.Printf("Error saving ledger: %v\n", err)
		return
	}
	file := filepath.Join(location, entry.filename)

	if inBrowser {
		if err := web.Open(file); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Printf("Opening up '%s' in Chrome...\n", entry.filename)
		return
	}

	fmt.Println("Don't forget to save (Ctrl + S)")
	cmd := exec.Command("open", "-a", "TextEdit", file)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}
