package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"path/filepath"

	// "sort"
	"slices"
)

var location string

// error checker
func check(e error) {
	if e != nil {
		panic(e)
	}
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

	// fmt.Println(os.Args)
	var filename string

	// check if file exists
	fileExists := func(filepath string) bool {
		_, err := os.Stat(filepath)
		return !os.IsNotExist(err)
	}

	// append to existing file
	appendToFile := func(filepath string, content string) {
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

	if len(os.Args) == 1 {
		filename = time.Now().Format("2006-01-02")
		filepath := location+"\\"+filename
		content := time.Now().Format("Monday, Jan 2, 2006")+"\n\n"
		
		if fileExists(filepath) {
			appendToFile(filepath, "\n\n" + content)
		} else {
			err := os.WriteFile(filepath, []byte(content), 0755)
			check(err)
		}

	} else if len(os.Args) > 2 && !slices.Contains(os.Args, "-t") {
		// fmt.Println(location)

		filename = strings.Join(os.Args[1:], " ")
		filepath := location+"\\"+filename
		content := time.Now().Format("Monday, Jan 2, 2006") + "\n\n"
		
		if fileExists(filepath) {
			appendToFile(filepath, "\n\n" + content)
		} else {
			err := os.WriteFile(filepath, []byte(content), 0755)
			check(err)
		}

	} else if slices.Contains(os.Args, "-t") {
		tIndex := findIndex(os.Args, "-t")
		// fmt.Printf("tIndex: %d", tIndex)
		if tIndex == -1 || tIndex == len(os.Args)-1 {
			fmt.Println("Invalid usage: '-t' flag requires a tag argument")
			return
		}

		if os.Args[1] == "-t" {
			filename = time.Now().Format("2006-01-02")
		} else {
			filename = strings.Join(os.Args[1:tIndex], " ")
		}

		tags := strings.Join(os.Args[tIndex+1:], ", ")
		filepath := location+"\\"+filename
		content := time.Now().Format("Monday, Jan 2, 2006") + "\n\n"

		if fileExists(filepath) {
			content = "\n\n" + "tags: " + tags + "\n" + content
			appendToFile(filepath, content)
		} else {
			err := os.WriteFile(filepath, []byte("tags: " + tags + "\n" + filename+ "\n"+ content +"\n"), 0755)
			check(err)
		}
	}

	fmt.Println("Don't forget to save (Ctrl + S)")
	// open the file in Notepad
	notepad_path := "C:\\Windows\\system32\\notepad.exe"
	file := fmt.Sprintf(location + "\\" + filename)
	cmd := exec.Command(notepad_path, file)
	err2 := cmd.Start() // Non-blocking program run
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
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
