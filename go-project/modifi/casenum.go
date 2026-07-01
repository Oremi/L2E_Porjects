package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {

	text := "this is a, sample , string for (up, 3)"

	words := strings.Fields(text)

	for i := 0; i < len(words); i++ {

		// for (up, N), (low, N) and (cap, N)
		if words[i] == "(up," || words[i] == "(low," || words[i] == "(cap," {

			// ensure the word exists
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
	fmt.Println(words)
}

func RemoveWords(words []string, index int) []string {
	result := append(words[:index], words[index+2:]...)
	return result

}
