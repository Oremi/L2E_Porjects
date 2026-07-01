package main

import (
	"fmt"
	"strings"
)

func isPunct(r rune) bool {
	return r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';'
}

func check(text string) string {

	runes := []rune(text)
	var build strings.Builder
	var last rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == ' ' && last == ' ' {
			continue
		}

		if isPunct(r) {
			if last == ' ' {
				out := build.String()
				build.Reset()
				build.WriteString(out[:len(out)-1])
			}
			for i < len(runes) && isPunct(runes[i]) {
				build.WriteRune(runes[i])
				last = runes[i]
				i++
			}

			if i < len(runes) && runes[i] != ' ' {
				build.WriteRune(' ')
				last = ' '
			}
			
			i--
			continue
		}

		build.WriteRune(r)
		last = r
	}

	return build.String()
}

func main() {
	text := "this is a string with a ' quote '    and ,there should be less space too !!"

	result := check(text)
	fmt.Println(text)
	fmt.Println()
	fmt.Println(result)
}
