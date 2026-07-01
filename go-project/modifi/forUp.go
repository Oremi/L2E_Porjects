package main

import (
	"fmt"
	"strconv"
	"strings"
)

func CaseConversion(words []string) []string {

	for i := 0; i < len(words); i++ {

		// strings.Tolower, Toupper
		if words[i] == "(up)" && i > 0 {
			words[i-1] = strings.ToUpper(words[i-1])
			words = RemoveWords(words, i)
			i--
		}

		// change to lowercase
		if words[i] == "(low)" && i > 0 {
			words[i-1] = strings.ToLower(words[i-1])
			words = RemoveWords(words, i)
			i--
		}

		// // for capitalize
		if words[i] == "(cap)" && i > 0 {
			words[i-1] = strings.ToUpper(words[i-1][0:1]) + strings.ToLower(words[i-1][1:])
			words = RemoveWords(words, i)
			i--
		}

		// for (up, N), (low, N) and (cap, N)
		if words[i] == "(up," || words[i] == "(low," || words[i] == "(cap," {

			//ensure that the words exist
			if i+1 >= len(words) {
				continue
			}
			ndx := strings.TrimSuffix(words[i+1], ")")
			ndxInt, err := strconv.Atoi(ndx)
			if err != nil {
				fmt.Printf("Error converting %s to integer: %v\n", ndx, err)
				continue
			}
			if i-ndxInt < 0 {
				fmt.Printf("Error: Not enough words to convert for %s\n", words[i])
				continue
			}

			for j := i - ndxInt; j < i; j++ {
				if words[i] == "(up," {
					words[j] = strings.ToUpper(words[j])
				} else if words[i] == "(low," {
					words[j] = strings.ToLower(words[j])
				} else if words[i] == "(cap," {
					words[j] = strings.ToUpper(words[j][0:1]) + strings.ToLower(words[j][1:])
				}
			}
			words = RemoveWords(words, i)
			i--
		}

	}
	return words
}
