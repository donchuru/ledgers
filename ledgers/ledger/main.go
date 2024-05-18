package main

import(
	"fmt"
	"os"
	"time"
	"strings"
	"bufio"
	"os/exec"
	// "sort"
	"slices"
)

func err_msg (e error) {
	if e != nil {
		fmt.Printf("Unable to write file: %v", e)
	}
}

func main () {
	/* take in command line arguments
	User inputs:
		ledger -> make a new ledger named today's date
		ledger "new Doc"  -> make a new ledger named new Doc
		ledger "new Doc" -t life -> make a new ledger named new Doc with the life tag
	*/

	// fetch location where all journals are stored
	f, _ := os.Open("../config/init.txt")
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	scanner.Scan()
	location := scanner.Text()

	// fmt.Println(os.Args)
	var filename string

	if len(os.Args) == 1 {
		filename = time.Now().Format("2006-01-02")
		err := os.WriteFile(location + "\\" + filename, []byte(time.Now().Format("Monday, Jan 2, 2006") + "\n"), 0755)
		err_msg(err)
	} else if len(os.Args) > 2 && !slices.Contains(os.Args, "-t") {
		// fmt.Println(location)

		filename = strings.Join(os.Args[1:], " ")
		err := os.WriteFile(location + "\\" + filename, []byte(filename + "\n" + time.Now().Format("Monday, Jan 2, 2006") + "\n\n"), 0755)
		err_msg(err)

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

		err := os.WriteFile(location + "\\" + filename, []byte("tags: " + tags + "\n" + filename + "\n" + time.Now().Format("Monday, Jan 2, 2006") + "\n\n"), 0755)
		err_msg(err)

	}

	// open the file
	exepath := "C:\\Windows\\system32\\notepad.exe"
	file := fmt.Sprintf(location + "\\" + filename)
	cmd := exec.Command(exepath, file)
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
