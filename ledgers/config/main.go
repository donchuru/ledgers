package main

import (
	"fmt"
	"os"
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

	fmt.Println("Prompt 1/2: Please enter your name: ")
	fmt.Scan(&name)

	fmt.Println("Prompt 2/2: Please enter where you'd like to store your journals (Paste with Ctrl+Shift+V): ")
	fmt.Scan(&location)

	d1 := []byte(name + "\n" + location)
    err := os.WriteFile("./init.txt", d1, 0644)
    check(err)
}