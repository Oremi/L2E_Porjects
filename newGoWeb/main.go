package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed banners/*.txt templates/* static/*
var content embed.FS

const charHeight = 8

var banners = make(map[string][]string)

type PageData struct {
	Title, Result string
}

func main() {
	// Load banners BEFORE starting the server/
	err := loadBanners()
	if err != nil {
		log.Fatalf("Failed to load banners: %v", err)
	}

	mux := http.NewServeMux()

	staticContent, err := fs.Sub(content, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))

	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/ascii-art", AsciiHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	render(w, "index", PageData{Title: "ASCII ART GENERATOR"})
}

func loadBanners() error {
	files := []string{"standard", "shadow", "thinkertoy"}

	for _, name := range files {
		content, err := content.ReadFile("banners/" + name + ".txt")
		if err != nil {
			return err
		}
		// Normalize line endings and split into a slice
		lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
		banners[name] = lines
	}
	return nil
}

func AsciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		errorWithRedirect(w, http.StatusMethodNotAllowed, "Method Not Allowed", 3, "/")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	bannerName := r.FormValue("banner")

	for _, r := range text {
		if (r < 32 || r > 126) && r != '\n' && r != '\r' {
			http.Error(w, "400 Bad Request: Non-ASCII characters detected", http.StatusBadRequest)
			return
		}
	}

	font, exists := banners[bannerName]
	if !exists {
		// http.Error(w, "400 Bad Request: Invalid banner", http.StatusBadRequest)
		errorWithRedirect(w, http.StatusBadRequest, "Invalid banner selected", 2, "/")
		return
	}

	inputLines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")

	var result strings.Builder
	for _, lines := range inputLines {
		if lines == "" {
			result.WriteString("\n")
			continue
		}
		result.WriteString(generateAscii(lines, font))
	}

	render(w, "index", PageData{Title: "ASCII ART GENERATOR", Result: result.String()})
}

func generateAscii(text string, font []string) string {
	var result strings.Builder

	for i := 1; i <= charHeight; i++ {
		for _, char := range text {
			index := (int(char)-32)*9 + i
			if index < len(font) {
				result.WriteString(font[index])
			}
		}
		result.WriteString("\n")
	}
	return result.String()
}

func render(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFS(content, "templates/header.html", "templates/"+tmpl+".html")
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute the specific page template
	// Use ExecuteTemplate to ensure you start with the main page, not the header
	err = t.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, "Execution Error: "+err.Error(), http.StatusInternalServerError)
	}
}

func errorWithRedirect(w http.ResponseWriter, status int, message string, delay int, redirectURL string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	htmlContent := fmt.Sprintf(` <DOCTYPE html>
<html>
<head>
	<meta http-equiv="refresh" content="%d;url=%s">
	<title>Error</title>
</head>
<body>
	<h1>Error %d</h1>
	<p>%s</p>
	<p>You will be redirected in %d seconds. If not, click <a href="%s">here</a>.</p>
</body>
</html>`, delay, redirectURL, status, message, delay, redirectURL)

	fmt.Fprint(w, htmlContent)
}
