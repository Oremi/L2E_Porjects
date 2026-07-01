
package main

import (
	"fmt"
)

func isPunc(r rune) bool {
	return r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';'
}

func main() {
	text := "this is a boy, whose name is unknown."
	//runes := []rune(text)
	
	found := false

	for i, r := range text {
		
		if isPunc(r) {
			fmt.Printf("found punctuation %v at position %v\n", string(r), i)
			found = true
		}

	}
	if !found {
		fmt.Println("No punctuations found in text")
	}	
}
