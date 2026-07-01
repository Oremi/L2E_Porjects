package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal()
	}

	fonts, err := os.ReadFile("standard.txt")
	// fonts, err := os.ReadFile("shadow.txt")
	// fonts, err := os.ReadFile("thinkertoy.txt")
	if err != nil {
		log.Fatal(err)
	}

	stringz := strings.Split(os.Args[1], "\\n")
	// fmt.Println(len(stringz))

	fontz := strings.ReplaceAll(string(fonts), "\r\n", "\n")
	splitLines := strings.Split(fontz, "\n")

	for i, chr := range stringz {
		if len(stringz) == 2 && i > 0 && stringz[i-1] == "" && chr == "" {

			break
		} else if chr == "" {

			fmt.Print("\n")
			continue
		}
		for i := 0; i < 8; i++ {

			for _, r := range chr {

				start := (int(r)-32)*9 + 1
				fmt.Print(splitLines[start+i])

			}
			fmt.Println()
		}
	}

}
