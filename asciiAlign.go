package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
)

//go:embed standard.txt shadow.txt thinkertoy.txt
var fonts embed.FS

func main() {
	var color, align string
	flag.StringVar(&color, "color", "", "Color name (red, green, blue, yellow)")
	flag.StringVar(&align, "align", "left", "left, center, right, justify")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: go run . [--color=red] [--align=center] [text] [banner]")
		return
	}

	// 1. Assign text and banner
	text := args[0]
	banner := "standard"
	if len(args) > 1 {
		banner = args[1]
	}

	// 2. Load Font Data
	data, err := fonts.ReadFile(banner + ".txt")
	if err != nil {
		fmt.Printf("Error: Banner '%s' not found\n", banner)
		return
	}
	fontLines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	// 3. Setup Variables
	ansiColor := getAnsi(color)
	termWidth := 80 // Simplified terminal width
	inputLines := strings.Split(strings.ReplaceAll(text, "\\n", "\n"), "\n")

	for _, line := range inputLines {
		if line == "" {
			fmt.Println()
			continue
		}

		// Calculate total width of the ASCII art for this specific line
		artWidth := 0
		spaceCount := 0
		for _, char := range line {
			if char == ' ' {
				spaceCount++
			}
			fIdx := int(char-32)*9 + 1
			artWidth += len(fontLines[fIdx])
		}

		// 4. Print the 8 rows of ASCII
		for row := 0; row < 8; row++ {
			
			// Handle Left/Center/Right Padding
			if align == "center" {
				fmt.Print(strings.Repeat(" ", max(0, (termWidth-artWidth)/2)))
			} else if align == "right" {
				fmt.Print(strings.Repeat(" ", max(0, termWidth-artWidth)))
			}

			// Handle Justify Variables
			extraPerSpace := 0
			remainder := 0
			if align == "justify" && spaceCount > 0 {
				totalExtra := termWidth - artWidth
				extraPerSpace = totalExtra / spaceCount
				remainder = totalExtra % spaceCount
			}

			// Print each character in the line
			for _, char := range line {
				fIdx := int(char-32)*9 + 1 + row
				
				if ansiColor != "" {
					fmt.Print(ansiColor)
				}
				fmt.Print(fontLines[fIdx])
				fmt.Print("\033[0m") // Reset color

				// If justifying, add extra spaces after the ASCII space character
				if char == ' ' && align == "justify" {
					add := extraPerSpace
					if remainder > 0 {
						add++
						remainder--
					}
					fmt.Print(strings.Repeat(" ", max(0, add)))
				}
			}
			fmt.Println()
		}
	}
}

func getAnsi(name string) string {
	colors := map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"blue":   "\033[34m",
		"yellow": "\033[33m",
	}
	return colors[strings.ToLower(name)]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

