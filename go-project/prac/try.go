package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--output" {
			fmt.Println("Usage: go run try.go --output=filename")
		}
	}
	outputFile := flag.String("output", "", "Name of the output file")
	flag.Parse()

	// args := flag.Args()
	// if len(args) < 1 {
	// 	fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
	// 	fmt.Println("Example: go run . --output=test.txt \"hello\" standard")
	// 	os.Exit(1)
	// }
	// Assign with defaults for the banner
	// input := args[0]
	// banner := "standard" // Default banner
	// if len(args) > 1 {
	// banner = args[1]
	// }

	if *outputFile == "" {
		fmt.Println("Please specify an output file with --output=filename")
		os.Exit(1)
	}

	if *outputFile != "" {
		fmt.Printf("Output will be written to: %s\n", *outputFile)

		result := strings.Builder{}
		result.WriteString("Hello, World!")
		result.WriteString("\n")

		os.WriteFile(*outputFile, []byte(result.String()), 0644)
	}
	// if *outputFile == "" {
	// 	fmt.Println("No output file specified. Output will be printed to console.")
	// 	fmt.Println("Hello, World!")
	// 	fmt.Printf("Printing ASCII art of '%s' using %s to console:\n", input, banner)
	// 	// Insert console printing logic here
	// }
}
