package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(Format("I was sitting over there ,and then BAMM !!"))
        fmt.Println(Format("As Elton John said: ' I am the most well-known homosexual '"))

}

func Format(s string) string {
	// 1. Move punctuation: "word ," -> "word,"
	for _, p := range []string{".", ",", "!", "?", ":", ";"} {
		s = strings.ReplaceAll(s, " "+p, p)   // Remove space before
		s = strings.ReplaceAll(s, p, p+" ")   // Ensure space after
	}

	// 2. Fix Quotes: "' word '" -> "'word'"
	parts := strings.Split(s, "'")
	for i := 1; i < len(parts); i += 2 {
		parts[i] = strings.TrimSpace(parts[i]) // Strip spaces inside quotes
	}
	s = strings.Join(parts, "'")

	// 3. Clean up: Fix "... " and double spaces
	s = strings.ReplaceAll(s, ". . .", "...")
	return strings.Join(strings.Fields(s), " ")
}

