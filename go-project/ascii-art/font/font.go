package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const charHeight = 8

var filesystem = os.DirFS(".")

func main() {

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run . \"text\" [fontFile]")
		fmt.Println("[fontFile]: Standard font: standard.txt, Shadow font: shadow.txt, Thinkertoy font: thinkertoy.txt")
		return
	}

	// input := os.Args[1]
	input := strings.ReplaceAll(os.Args[1], "\\n", "\n")
	lines := strings.Split(input, "\n")

	//choosing the file style
	fontFile := "standard.txt"
	if len(os.Args) == 3 {
		fontFile = os.Args[2]
	}

	fontLines, err := loadFont(fontFile)
	if err != nil {
		fmt.Println("Error reading font:", err)
		fmt.Println("Falling back to standard font")

		fontLines, err = loadFont("standard.txt")
		if err != nil {
			fmt.Println("Error reading font:", err)
			return
		}
	}

	fontMap := buildFontMap(fontLines)
	for _, line := range lines {
		renderText(line, fontMap)
	}
}

func loadFont(filename string) ([]string, error) {

	data, err := fs.ReadFile(filesystem, filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func buildFontMap(lines []string) map[rune][]string {

	font := make(map[rune][]string)

	ascii := 32
	i := 0

	for ascii <= 126 {

		font[rune(ascii)] = lines[i : i+charHeight]

		i += charHeight + 1
		ascii++
	}

	return font
}

func renderText(text string, font map[rune][]string) {

	for row := 0; row < charHeight; row++ {

		for _, char := range text {

			if letter, ok := font[char]; ok {
				fmt.Print(letter[row])
			}
		}

		fmt.Println()
	}
}
