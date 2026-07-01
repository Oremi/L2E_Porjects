package main

import (
	"fmt"
	"strings"
)

func fixArticles(text string) string {
	words := strings.Fields(text)
	vowelsAndH := "aeiouhAEIOUH"

	for i := 0; i < len(words)-1; i++ {
		stripped := words[i]
		if strings.ToLower(stripped) == "a" {
			nextWord := words[i+1]
			firstLetter := string(nextWord[0])

			if strings.ContainsAny(firstLetter, vowelsAndH) {
				words[i] = "an"
			}
		}
	}

	return strings.Join(words, " ")
}

func main() {
	text := "this is a egg string and a honest mistake"
	fmt.Println(fixArticles(text))
}
