package main

import(
	"fmt"
	"os"
	"time"
	"strings"
	"bufio"
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
	*/
	if len(os.Args) == 1 {
		today := time.Now().Format("2006-01-02")
		err := os.WriteFile(today, []byte(time.Now().Format("Monday, Jan 2, 2006") + "\n"), 0755)
		err_msg(err)
	} else if len(os.Args) > 2 {
		
		f, _ := os.Open("../config/init.txt")
		scanner := bufio.NewScanner(f)
		scanner.Scan()
		scanner.Scan()
		location := scanner.Text()

		// fmt.Println(location)

		err := os.WriteFile(location + "\\" + strings.Join(os.Args[1:], " "), []byte(strings.Join(os.Args[1:], " ") + "\n" + time.Now().Format("Monday, Jan 2, 2006") + "\n\n"), 0755)
		err_msg(err)
	}
}
