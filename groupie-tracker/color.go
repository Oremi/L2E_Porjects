package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

//go:embed standard.txt
var standardFont string

// colorFlags is a custom type to allow the 'flag' package to handle multiple --color inputs.
type colorFlags []string

func (i *colorFlags) String() string     { return strings.Join(*i, ",") }
func (i *colorFlags) Set(v string) error { *i = append(*i, v); return nil }

const charHeight = 8

func main() {
	// Reject "--color red" format
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--color" {
			fmt.Println("Usage: go run . [OPTION] [STRING]")
			fmt.Println()
			fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
			return
		}
	}

	var colors colorFlags
	flag.Var(&colors, "color", "Color (name, #hex, rgb(r,g,b), or hsl(h,s,l))")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: go run . --color=red [substring] \"text\"")
		return
	}

	// Determine if the user provided a specific substring to colorize.
	// If only one argument is provided, it's treated as the full text.
	text, sub := args[len(args)-1], ""
	if len(args) > 1 {
		sub = args[0]
	}

	// We pass os.Stdout to the renderer
	if err := RenderASCII(os.Stdout, text, sub, colors); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// RenderASCII processes the input text and writes the ASCII art representation to the provided io.Writer.
func RenderASCII(w io.Writer, text, sub string, colorInputs []string) error {
	font := strings.Split(strings.ReplaceAll(standardFont, "\r\n", "\n"), "\n")
	if len(font) < 2 {
		return fmt.Errorf("font file is empty or invalid")
	}

	// Pre-parse all input colors into ANSI escape sequences.
	var ansiCodes []string
	for _, c := range colorInputs {
		if code := parseColor(c); code != "" {
			ansiCodes = append(ansiCodes, code)
		}
	}

	// Split input text by literal newline characters.
	lines := strings.Split(strings.ReplaceAll(text, "\\n", "\n"), "\n")
	for _, line := range lines {
		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		// Each ASCII character is charHeight rows tall; we must print row-by-row
		for row := 0; row < charHeight; row++ {
			colorIdx := 0
			for i := 0; i < len(line); i++ {
				// Check if this specific character should be colored based on substring match.
				isColored := len(ansiCodes) > 0 && (sub == "" || isInside(line, sub, i))

				if isColored {
					// Cycle through provided colors if multiple were passed.
					fmt.Fprint(w, ansiCodes[colorIdx%len(ansiCodes)])
					colorIdx++
				}

				// Calculate index: (ASCII value - 32) * (height + 1 separator line) + 1 offset + current row.
				fIdx := int(line[i]-32)*9 + 1 + row
				if fIdx < len(font) {
					fmt.Fprint(w, font[fIdx])
				}
				if isColored {
					fmt.Fprint(w, "\033[0m") // Reset color after printing the character slice.
				}
			}
			fmt.Fprintln(w)
		}
	}
	return nil
}

// parseColor converts various color string formats into ANSI truecolor (24-bit) or 8-bit escape sequences.
func parseColor(c string) string {
	c = strings.ToLower(strings.TrimSpace(c))
	var r, g, b int

	// Handle standard terminal colors and a custom 'orange'.
	names := map[string]string{"red": "31", "green": "32", "yellow": "33", "blue": "34", "magenta": "35", "cyan": "36", "orange": "38;5;208"}
	if code, ok := names[c]; ok {
		return "\033[" + code + "m"
	}

	// HEX format (e.g., #FF0000)
	if _, err := fmt.Sscanf(c, "#%02x%02x%02x", &r, &g, &b); err == nil {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	// RGB format (e.g., rgb(255,0,0))
	if _, err := fmt.Sscanf(c, "rgb(%d,%d,%d)", &r, &g, &b); err == nil {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	// HSL format (e.g., hsl(0,100,50))
	var h, s, l float64
	if _, err := fmt.Sscanf(strings.ReplaceAll(c, "%", ""), "hsl(%f,%f,%f)", &h, &s, &l); err == nil {
		r, g, b := hslToRgb(h/360, s/100, l/100)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	return ""
}

// hslToRgb converts HSL values to RGB (0-255).
// Implementation uses the formula for mapping hue to color cylinders.
func hslToRgb(h, s, l float64) (int, int, int) {
	f := func(n float64) float64 {
		k := math.Mod(n+h*12, 12)
		a := s * math.Min(l, 1-l)
		return l - a*math.Max(-1, math.Min(math.Min(k-3, 9-k), 1))
	}
	return int(f(0) * 255), int(f(8) * 255), int(f(4) * 255)
}

// isInside returns true if the character index 'i' in the string 'text'
// falls within any occurrence of the substring 'sub'.
func isInside(text, sub string, i int) bool {
	for start := 0; ; {
		idx := strings.Index(text[start:], sub)
		if idx == -1 {
			break
		}
		if i >= start+idx && i < start+idx+len(sub) {
			return true
		}
		start += idx + 1 // Advance search to catch overlapping or multiple occurrences.
	}
	return false
}
