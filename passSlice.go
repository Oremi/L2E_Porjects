package main

import (
	"math/rand/v2"
	"fmt"
)


func main() {

	letters := Alphabet()
	symbols := Symbol()
	numbers := Number()

	var letterCount, symbolCount, numberCount int
	
	fmt.Println(`
	******************************************************************

			WELCOME TO THE PASSWORD GENERATOR


	******************************************************************`)

	fmt.Print("How many alphabets do you want?: ")
	fmt.Scan(&letterCount)

	fmt.Print("How many symbols do you want?: ")
        fmt.Scan(&symbolCount)

	fmt.Print("How many numbers do you want?: ")
        fmt.Scan(&numberCount)

	totalCount := letterCount + symbolCount + numberCount

	password := make([]string, totalCount)

	for i := 0; i < letterCount; i++ {
		password[i] = letters[rand.IntN(len(letters))]
	}

	for i := letterCount; i < letterCount + symbolCount; i++ {
                password[i] = symbols[rand.IntN(len(symbols))]
        }

	for i := letterCount + symbolCount; i < totalCount; i++ {
                password[i] = numbers[rand.IntN(len(numbers))]
        }

	rand.Shuffle(totalCount, func(i, j int){
		password[i], password[j] = password[j], password[i]
	})

	passwordString := join(password)

	fmt.Printf("This is your password: %v \n", passwordString)

}

func Alphabet() []string {
	var num int = 0
	letters := make([]string, 26)
	for i := 0; i < 26; i++ {
		letters[i] = string(rune('a' + num))
		num++
	}
	return letters
}

func Symbol() []string {
        var num int = 0
	symbols := make([]string, 15)
        for i := 0; i < 15; i++ {
                symbols[i] = string(rune('!' + num))
		num++
        }
        return symbols
}

func Number() []string {
	var num int = 0
        numbers := make([]string, 10)
        for i := 0; i < 10; i++ {
                numbers[i] = string(rune('0' + num))
        	num++
	}
        return numbers
}

func join(slice []string) string {
	var letters string
	for _, char := range slice {
		letters += char
	}
	return letters
}


