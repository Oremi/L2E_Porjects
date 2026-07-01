package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	// getting a user input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	user_input := scanner.Text()

	fmt.Printf("%s \n", user_input)

	var s string = "8.2"
	var b string = "10"

	a1, _ := strconv.ParseFloat(s, 64)
	b2, _ := strconv.Atoi(b)

	c := a1 + float64(b2)

	fmt.Printf("%f \n", c)
}
