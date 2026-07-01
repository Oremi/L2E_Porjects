package main

import (
	"fmt"
	"learn"
	"strings"
)

func main() {

	text := "this is a boy 1E (hex) and this is 10 (bin)"
	new_text := strings.Fields(text)

	result := learn.DecimalConversion(new_text)

	fmt.Println(result[:])

	// for i := 0; i < len(new_text); i++ {

	// 	if new_text[i] == "(hex)" && i > 0 {
	// 		if value, err := strconv.ParseInt(new_text[i-1], 16, 64); err == nil {
	// 			new_text[i-1] = strconv.FormatInt(value, 10)
	// 			new_text = append(new_text[:i], new_text[i+1:]...)
	// 			i--
	// 		}
	// 		fmt.Println(new_text[:])

	// 		// fmt.Printf("this is the type %T", new_text[i-1])
	// 	}

	// 	if new_text[i] == "(bin)" && i > 0 {
	// 		if value, err := strconv.ParseInt(new_text[i-1], 2, 64); err == nil {
	// 			new_text[i-1] = strconv.FormatInt(value, 10)
	// 			new_text = append(new_text[:i], new_text[i+1:]...)
	// 		}
	// 		fmt.Println(new_text[:])

	// 		// fmt.Printf("this is the type %T", new_text[i-1])
	// 	}

	// }

	// fmt.Println(new_text[:])
	// fmt.Println(new_text[1:])

}
