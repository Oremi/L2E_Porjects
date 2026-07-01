package main

import (
	"fmt"
	"strings"
	"regexp"
)

func main() {
	text := "this is a boy to match (bin)"	
	words := strings.Fields(text)
	markerPattern := regexp.MustCompile('\((hex|bin)\)$')

	for i := 0; i < len(words); i++ {

	}
}

func PunctuationHandler(text string) string {
	// remove space before punctuation
	re1 := regexp.MustCompile(`\s+([.,!?;:])`)
	text = re1.ReplaceAllString(text, "$1")

	// ensure single space after punctuation
	re2 := regexp.MustCompile(`([.,!?;:])([^\s])`)
	text = re2.ReplaceAllString(text, "$1 $2")

	// remove multiple spaces
	re3 := regexp.MustCompile(`\s+`)
	text = re3.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

func QuoteHandler(text string) string {
	
	re := regexp.MustCompile(`'\s*(.*?)\s*'`)
	s = re.ReplaceAllString(s, "'$1'")

	re2 := regexp.MustCompile(`"\s*(.*?)\s*"`)
	s = re2.ReplaceAllString(s, "\"$1\"")

	return strings.TrimSpace(text)
}
