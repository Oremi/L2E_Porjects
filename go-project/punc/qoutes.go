package main

import (
	"fmt"
	"strings"
)

func IsVowelOrH(c rune) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'h':
		return true
	}
	return false
}

func Qoutes(words []string) []string {
	
	for i, word := range words {
		// Check if the current word is "a" or "A"
		if (word == "a" || word == "A") && i+1 < len(words) {
			// Look at the first character of the next word
			nextWord := words[i+1]
			firstRune := rune(nextWord[0])

			if IsVowelOrH(firstRune) {
				if word == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}
	//return strings.Join(words, " ")
	return words
}
func main() {
	text := "A egg for this is a sample text to text vowels, a hungry man"
	words := strings.Fields(text)
	words = Qoutes(words)
	fmt.Println(words)
}

