package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestLoadFont(t *testing.T) {
	lines, err := loadFont("standard.txt")

	if err != nil {
		t.Fatalf("expected font to load, got error: %v", err)
	}

	if len(lines) == 0 {
		t.Fatal("font file should not be empty")
	}
}

func TestBuildFontMap(t *testing.T) {
	lines, err := loadFont("standard.txt")
	if err != nil {
		t.Fatal(err)
	}

	font := buildFontMap(lines)

	if len(font) != 95 {
		t.Errorf("expected 95 ascii characters, got %d", len(font))
	}

	if _, ok := font['A']; !ok {
		t.Error("expected font map to contain 'A'")
	}

	if _, ok := font[' ']; !ok {
		t.Error("expected font map to contain space")
	}
}

func TestRenderText(t *testing.T) {
	lines, err := loadFont("standard.txt")
	if err != nil {
		t.Fatal(err)
	}

	font := buildFontMap(lines)

	output := captureOutput(func() {
		renderText("A", font)
		renderText("Hello World", font)
	})

	if strings.TrimSpace(output) == "" {
		t.Error("expected ascii art output for 'A'")
	}
}

func TestRenderText_UnknownCharacter(t *testing.T) {
	lines, err := loadFont("standard.txt")
	if err != nil {
		t.Fatal(err)
	}

	font := buildFontMap(lines)

	output := captureOutput(func() {
		renderText("\t", font)
	})

	if output == "" {
		t.Error("expected rows even if character is unknown")
	}
}
