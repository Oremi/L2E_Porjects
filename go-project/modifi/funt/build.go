package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "this is a , sample string for you . hello .... biy "
	result := FixP(text)
	fmt.Println(result)
}

func FixP(text string) string {

	var result strings.Builder

	for i := 0; i < len(text); i++ {
		char := text[i]

		if strings.ContainsRune(",.!?;:", rune(char)) {
			if result.Len() > 0 {
				prev := result.String()[result.Len()-1]
				if prev == ' ' {
					resultStr := result.String()
					result.Reset()
					result.WriteString(resultStr[:len(resultStr)-1])
				}
			}

			result.WriteByte(char)

			if i+1 < len(text) && !strings.ContainsRune(",.!?;: ", rune(text[i+1])) && text[i+1] != ' ' {
				result.WriteByte(' ')
			}

		} else {
			result.WriteByte(char)
		}
	}
	return result.String()
}
