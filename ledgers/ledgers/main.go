package main

import(
	"fmt"
	"os"
	"log"
)

func main () {
	/* take in command line arguments
	User inputs:
		ledgers -> make a new ledger named today's date
		ledger "new Doc"  -> make a new ledger named new Doc
	*/
	entries, err := os.ReadDir("../your_journals")
	if len(os.Args) == 1 {
		// show all journals
		if err != nil {
			log.Fatal(err)
		}
	
		for _, e := range entries {
			fmt.Println(e.Name())
		}

	} else if len(os.Args) == 2 {
		if os.Args[1] == "-m"{
			// show me list of all journals in order of last modified
			for _, e := range entries {
				fmt.Println(e.Name())
			}

		} else if os.Args[1] == "-c" {
			// show me list of all journals in order of last created
			for _, e := range entries {
				fmt.Println(e.Name())
			}
		}
	}
}




/* for ledgers:

*/