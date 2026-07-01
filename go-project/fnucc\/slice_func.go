package main

import (
	"strconv"
)

func DecimalConversion(words []string) []string {

	for i := 0; i < len(words); i++ {

		// for hexadecimal conversion
		if words[i] == "(hex)" && i > 0 {
			if val, err := strconv.ParseInt(words[i-1], 16, 64); err == nil {
				words[i-1] = strconv.FormatInt(val, 10)
				words = append(words[:i], words[i+1:]...)
				i--
			}

		}

		// for binary conversion
		if words[i] == "(bin)" && i > 0 {
			if val, err := strconv.ParseInt(words[i-1], 2, 64); err == nil {
				words[i-1] = strconv.FormatInt(val, 10)
				words = append(words[:i], words[i+1:]...)
				i--
			}

		}
	}
	return words

}
