package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

//go:embed standard.txt shadow.txt thinkertoy.txt
var fonts embed.FS

type colorFlags []string

func (i *colorFlags) String() string { return strings.Join(*i, ",") }
func (i *colorFlags) Set(v string) error {
	*i = append(*i, v)
	return nil
}

const charHeight = 8

func main() {
	var colors colorFlags
	var align string

	flag.Var(&colors, "color", "Color (name, #hex, rgb(r,g,b), or hsl(h,s,l))")
	flag.StringVar(&align, "align", "left", "Text alignment: left, center, right, justify")

	flag.Parse()
	args := flag.Args()

	// Validate align flag strictly
	validAlign := map[string]bool{
		"left": true, "center": true, "right": true, "justify": true,
	}
	if !validAlign[align] {
		printUsage()
		return
	}

	// Validate arguments count
	if len(args) == 0 || len(args) > 2 {
		printUsage()
		return
	}

	text := args[0]
	banner := "standard"

	if len(args) == 2 {
		banner = args[1]
	}
	
	fontData, err := loadFont(banner)
	if err != nil {
		printUsage()
		return
	}

	if err := RenderASCII(os.Stdout, text, "", colors, align, fontData); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: go run . [OPTION] [STRING] [BANNER]
Example: go run . --align=right something standard`)
}

func RenderASCII(w io.Writer, text, sub string, colorInputs []string, align string, fontData string) error {
	font := strings.Split(strings.ReplaceAll(fontData, "\r\n", "\n"), "\n")
	if len(font) < 2 {
		return fmt.Errorf("font file is empty or invalid")
	}

	// Get terminal width dynamically
	termWidth := getTerminalWidth()
	//if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
	//	termWidth = width
	//}

	var ansiCodes []string
	for _, c := range colorInputs {
		if code := parseColor(c); code != "" {
			ansiCodes = append(ansiCodes, code)
		}
	}

	lines := strings.Split(strings.ReplaceAll(text, "\\n", "\n"), "\n")

	for _, line := range lines {
		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		words := strings.Fields(line)

		// Precompute character widths
		charWidths := make([]int, len(line))
		totalCharsWidth := 0

		for i := 0; i < len(line); i++ {
			if line[i] == ' ' {
				charWidths[i] = 1
				totalCharsWidth += 1
				continue
			}
			fIdx := int(line[i]-32)*9 + 1
			wid := len(font[fIdx])
			charWidths[i] = wid
			totalCharsWidth += wid
		}

		for row := 0; row < charHeight; row++ {

			// Calculate actual row width
			actualLineWidth := 0
			for i := 0; i < len(line); i++ {
				if line[i] == ' ' {
					actualLineWidth += 1
					continue
				}
				fIdx := int(line[i]-32)*9 + 1 + row
				actualLineWidth += len(font[fIdx])
			}

			// Skip if overflow
			if actualLineWidth > termWidth {
				continue
			}

			switch align {

			case "center", "right":
				pad := termWidth - actualLineWidth
				if align == "center" {
					pad /= 2
				}
				if pad > 0 {
					fmt.Fprint(w, strings.Repeat(" ", pad))
				}
				printLineStandard(w, line, sub, font, ansiCodes, row)

			case "justify":
				if len(words) <= 1 {
					printLineStandard(w, line, sub, font, ansiCodes, row)
					break
				}

				totalSpaceNeeded := termWidth - totalCharsWidth
				if totalSpaceNeeded < 0 {
					printLineStandard(w, line, sub, font, ansiCodes, row)
					break
				}

				spaceCount := len(words) - 1
				baseSpace := totalSpaceNeeded / spaceCount
				extraSpace := totalSpaceNeeded % spaceCount

				wordIndex := 0
				i := 0

				for i < len(line) {
					if line[i] == ' ' {
						i++
						continue
					}

					// extract word
					start := i
					for i < len(line) && line[i] != ' ' {
						i++
					}
					word := line[start:i]

					printWord(w, word, line, sub, font, ansiCodes, row, start)

					if wordIndex < spaceCount {
						spaceSize := baseSpace
						if wordIndex < extraSpace {
							spaceSize++
						}
						fmt.Fprint(w, strings.Repeat(" ", spaceSize))
					}

					wordIndex++
				}

			default: // left
				printLineStandard(w, line, sub, font, ansiCodes, row)
			}

			fmt.Fprintln(w)
		}
	}

	return nil
}

func printLineStandard(w io.Writer, line, sub string, font []string, codes []string, row int) {
	colorIdx := 0

	for i := 0; i < len(line); i++ {
		isColored := len(codes) > 0 && (sub == "" || isInside(line, sub, i))

		if line[i] == ' ' {
			fmt.Fprint(w, " ")
			continue
		}

		fIdx := int(line[i]-32)*9 + 1 + row

		if isColored {
			fmt.Fprint(w, codes[colorIdx%len(codes)])
			colorIdx++
		}

		fmt.Fprint(w, font[fIdx])

		if isColored {
			fmt.Fprint(w, "\033[0m")
		}
	}
}

func printWord(w io.Writer, word, fullLine, sub string, font []string, codes []string, row int, startIdx int) {
	for i := 0; i < len(word); i++ {
		charIdx := startIdx + i
		isColored := len(codes) > 0 && (sub == "" || isInside(fullLine, sub, charIdx))

		fIdx := int(word[i]-32)*9 + 1 + row

		if isColored {
			fmt.Fprint(w, codes[0])
		}

		fmt.Fprint(w, font[fIdx])

		if isColored {
			fmt.Fprint(w, "\033[0m")
		}
	}
}

func parseColor(c string) string {
	c = strings.ToLower(strings.TrimSpace(c))

	var r, g, b int

	names := map[string]string{
		"red": "31", "green": "32", "yellow": "33",
		"blue": "34", "magenta": "35", "cyan": "36",
		"orange": "38;5;208",
	}

	if code, ok := names[c]; ok {
		return "\033[" + code + "m"
	}

	if _, err := fmt.Sscanf(c, "#%02x%02x%02x", &r, &g, &b); err == nil {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}

	if _, err := fmt.Sscanf(c, "rgb(%d,%d,%d)", &r, &g, &b); err == nil {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}

	var h, s, l float64
	if _, err := fmt.Sscanf(strings.ReplaceAll(c, "%", ""), "hsl(%f,%f,%f)", &h, &s, &l); err == nil {
		r, g, b := hslToRgb(h/360, s/100, l/100)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}

	return ""
}

func hslToRgb(h, s, l float64) (int, int, int) {
	f := func(n float64) float64 {
		k := math.Mod(n+h*12, 12)
		a := s * math.Min(l, 1-l)
		return l - a*math.Max(-1, math.Min(math.Min(k-3, 9-k), 1))
	}
	return int(f(0) * 255), int(f(8) * 255), int(f(4) * 255)
}

func isInside(text, sub string, i int) bool {
	for start := 0; ; {
		idx := strings.Index(text[start:], sub)
		if idx == -1 {
			break
		}
		if i >= start+idx && i < start+idx+len(sub) {
			return true
		}
		start += idx + 1
	}
	return false
}

func getTerminalWidth() int {
	ws := &winsize{}

	retCode, _, _ := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(retCode) == -1 || ws.Col == 0 {
		return 80 // fallback
	}

	return int(ws.Col)
}

func loadFont(name string) (string, error) {
	switch name {
	case "standard", "shadow", "thinkertoy":
		data, err := fonts.ReadFile(name + ".txt")
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("invalid banner")
	}
}
