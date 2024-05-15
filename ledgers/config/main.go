package main

import (
	"fmt"
	"os"
	"bufio"
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

	fmt.Println("Prompt 2/2: Please enter where you'd like to store your journals (Paste with Ctrl+Shift+V): ")
	if scanner.Scan() {
		location = scanner.Text()
		// fmt.Printf("Input was: %q\n", location)
	}

	d1 := []byte(name + "\n" + location)
    err := os.WriteFile("./init.txt", d1, 0644)
    check(err)
}