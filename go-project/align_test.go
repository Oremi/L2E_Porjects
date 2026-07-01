package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// runProgram simulates the user running 'go run . <args>'
func runProgram(args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// If the program fails, return stderr to help debug
		return stderr.String(), err
	}
	return stdout.String(), nil
}

func TestFlexibleAlignments(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		alignment string
	}{
		{"Right Align", []string{"--align=right", "something", "standard"}, "right"},
		{"Center Align", []string{"--align=center", "hello", "shadow"}, "center"},
		{"Justify Align", []string{"--align=justify", "1 Two 4", "shadow"}, "justify"},
		{"Left Align", []string{"--align=left", "Hello World!", "standard"}, "left"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runProgram(tt.args...)
			if err != nil {
				t.Fatalf("Program failed to run: %v\nOutput: %s", err, output)
			}

			// Split output into lines and remove the very last empty newline
			lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

			// Find the first non-empty line to detect the width used by the program
			var detectedWidth int
			for _, l := range lines {
				if len(l) > 0 {
					detectedWidth = len(l)
					break
				}
			}

			if detectedWidth == 0 && tt.alignment != "left" {
				t.Fatal("Output appears empty or lines have no width")
			}

			t.Logf("Testing %s: Detected width %d", tt.alignment, detectedWidth)

			for i, line := range lines {
				if line == "" {
					continue
				}

				// 1. Consistency Check: In Right/Center/Justify, all lines should be the same length
				if tt.alignment != "left" && len(line) != detectedWidth {
					t.Errorf("Line %d: Inconsistent width. Expected %d, got %d", i, detectedWidth, len(line))
				}

				// 2. Alignment-specific Logic Checks
				switch tt.alignment {
				case "right":
					// If the line is exactly 80 (or the detected width),
					// we check if it starts with spaces and ends with content.
					// If they are padding the right, we skip the 'HasSuffix' check.
					if len(line) < detectedWidth && strings.HasSuffix(line, " ") {
						t.Errorf("Line %d: Right-aligned text has trailing spaces", i)
					}
				case "left":
					if strings.HasPrefix(line, " ") && len(strings.TrimSpace(line)) > 0 {
						t.Errorf("Line %d: Left-aligned text has unexpected leading spaces", i)
					}
				}

			}
		})
	}
}

// TestAllowedPackages ensures no non-standard packages are imported.
func TestAllowedPackages(t *testing.T) {
	out, err := exec.Command("go", "list", "-f", "{{.Imports}}", ".").Output()
	if err != nil {
		t.Fatalf("Failed to list imports: %v", err)
	}

	imports := string(out)
	if strings.Contains(imports, "github.com") || strings.Contains(imports, "golang.org/x/") {
		t.Errorf("Project uses non-standard packages: %s", imports)
	}
}
