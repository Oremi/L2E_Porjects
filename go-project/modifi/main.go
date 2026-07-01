package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {

	text := "this is a string that 1E (hex) will become a slice , 1a (hex) note the 2E (hex) punctuation ."
	new_text := strings.Fields(text)

	result := DecimalConversion(new_text)
	result = CaseConversion(result)

	fmt.Println(result[:])
}

func DecimalConversion(words []string) []string {

	for i := 0; i < len(words); i++ {

		// for hexadecimal conversion
		if words[i] == "(hex)" && i > 0 {
			if val, err := strconv.ParseInt(words[i-1], 16, 64); err == nil {

				words[i-1] = strconv.FormatInt(val, 10)
				// Remove the "(hex)" from the slice and adjust the index (before the index of "(hex)" and after the index of "(hex)")
			}
			words = append(words[:i], words[i+1:]...)
			i--

		}

		// for binary conversion
		if words[i] == "(bin)" && i > 0 {
			if val, err := strconv.ParseInt(words[i-1], 2, 64); err == nil {

				words[i-1] = strconv.FormatInt(val, 10)
				// Remove the "(bin)" from the slice and adjust the index (before the index of "(bin)" and after the index of "(bin)")
			}
			words = append(words[:i], words[i+1:]...)
			i--

		}
	}
	return words

}

func CaseConversion(words []string) []string {

	for i := 0; i < len(words); i++ {

		// strings.Tolower, Toupper
		if words[i] == "(up)" && i > 0 {
			words[i-1] = strings.ToUpper(words[i-1])
			words = append(words[:i], words[i+1:]...)
			i--
		}

		// change to lowercase
		if words[i] == "(low)" && i > 0 {
			words[i-1] = strings.ToLower(words[i-1])
			words = append(words[:i], words[i+1:]...)
			i--
		}

		// // for capitalize
		if words[i] == "(cap)" && i > 0 {
			words[i-1] = strings.ToUpper(words[i-1][0:1]) + strings.ToLower(words[i-1][1:])
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
	return words
}
