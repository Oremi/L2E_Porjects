package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const charHeight = 8

func main() {
	if (len(os.Args) < 2) || (len(os.Args) > 3) {
		fmt.Println("Usage: go run . \"text\" [fontFile]")
		fmt.Println("[fontFile]: Standard font: standard.txt, Shadow font: shadow.txt, Thinkertoy font: thinkertoy.txt")
		return
	}

	// input := os.Args[1]
	lines := strings.ReplaceAll(os.Args[1], "\\n", "\n")
	input := strings.Split(lines, "\n")

	filesystem := os.DirFS(".")

	// choose font
	fontFile := "standard.txt"

	if len(os.Args) == 3 {
		fontFile = os.Args[2]
	}

	data, err := fs.ReadFile(filesystem, fontFile)
	if err != nil {
		fmt.Println("Error reading font:", err)
		fmt.Print("Falling back to standard font:")
		data, err = fs.ReadFile(filesystem, "standard.txt")
		if err != nil {
			fmt.Println("Error reading font:", err)
			return
		}
	}

	liness := strings.Split(string(data), "\n")

	// 1. Check if the input is only newlines (or empty)
	onlyNewlines := true
	for _, word := range input {
		if word != "" {
			onlyNewlines = false
			break
		}
	}

	// 2. Handle the output
	if onlyNewlines {
		// If the input was "\n", strings.Split makes ["", ""].
		// We print N-1 newlines to match standard behavior.
		if len(input) > 0 {
			for i := 0; i < len(input)-1; i++ {
				fmt.Println()
			}
		}
	} else {
		// Standard printing logic
		for _, line := range input {
			if line == "" {
				fmt.Println()
				continue
			}
			printAscii(line, liness)
		}
	}
}

func printAscii(text string, font []string) {

	for row := 0; row < charHeight; row++ {

		for _, char := range text {

			index := (int(char) - 32) * 9

			fmt.Print(font[index+row])
		}

		fmt.Println()
	}
}
