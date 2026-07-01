package main

import (
	"fmt"
	"strings"
	"strconv"
)

var demo_strings string

func main() {
	
	demo_strings = "what is 1E (hex) and 10 (bin)"

	words := strings.Fields(demo_strings)
	
	words = modif(words)
	
	fmt.Println("This is the original string of text %s", demo_strings)
	fmt.Println("This is the modified/converted string of texr %s", words)
}

func modif(words []string) []string {

	for i := 0; i < len(words); i++ {

		if words[i] == "(hex)" && i > 0 {
			if value, err := strconv.ParseInt(words[i-1], 16, 64); err == nil {
				words[i-1] = strconv.FormatInt(value, 10)
			}
			words = append(words[:i], words[i+1:]...)
			i--
		}

		if words[i] == "(bin)" && i > 0 {
			if value, err := strconv.ParseInt(words[i-1], 2, 64); err == nil {
				words[i-1] = strconv.FormatInt(value, 10)
			}
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
	return words
				
}	

