package main

import (
	_ "fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const charHeight = 8

type PageData struct {
	Result string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/ascii-art", asciiHandler)

	log.Println("Server starting at http://localhost:8080")
	// Removed 'go' so the main goroutine blocks and keeps the server running
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	renderTemplate(w, "index", nil)
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	bannerName := r.FormValue("banner")

	// 1. Validate and Read the font file
	filePath := "banners/" + bannerName + ".txt"
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "404 Banner Not Found", http.StatusNotFound)
		return
	}

	// 2. Prepare font slice (handling Windows/Unix line endings)
	font := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")

	// 3. Generate the art
	// We handle the split input (newlines) here
	inputLines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	var finalResult strings.Builder

	for _, line := range inputLines {
		if line == "" {
			finalResult.WriteString("\n")
			continue
		}
		finalResult.WriteString(generateAscii(line, font))
	}

	// 4. Render once
	renderTemplate(w, "index", PageData{Result: finalResult.String()})
}

func generateAscii(text string, font []string) string {
	var result strings.Builder

	for row := 0; row < charHeight; row++ {
		for _, char := range text {
			if char < 32 || char > 126 {
				continue
			}

			// Standard ASCII banner format: 9 lines per char (1 empty + 8 of art)
			index := 9*(int(char)-32) + 1

			if index+row < len(font) {
				result.WriteString(font[index+row])
			}
		}
		result.WriteString("\n")
	}
	return result.String()
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
