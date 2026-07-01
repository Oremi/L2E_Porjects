package main


import (
	"net/http"
	"html/template"
	"fmt"
	"os"
	"strings"
)

type Pagedata struct {
	Title	string
	Result	string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/asciiHandler", asciiHandler)
	fmt.Println("Server started at http://localhost:8000")

	http.ListenAndServe(":8000", mux)
}


func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		
		if r.URL.Path != "/" {
			http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
			return
		}

		renderTemplate(w, "index.html", Pagedata{
			Title: "ASCII ART GENERATOR",
		})
	}
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	userInput := r.FormValue("text")
	bannerName := r.FormValue("banner") //standard

	banner, err := os.ReadFile(bannerName + ".txt") //standard.txt
	if err != nil {
		http.Error(w, "400 Bad Request: Invalid banner", http.StatusBadRequest)
		return
	}

	for _, char := range userInput {
		if char < 32 || char > 126 {
			http.Error(w, "400 Bad Request: Non-ASCII characters detected", http.StatusBadRequest)
			return
		}
	}

	// format the banner

	font := strings.ReplaceAll(string(banner), "\r\n", "\n")
	realFont := strings.Split(font, "\n") 

	// format the text

	input := strings.ReplaceAll(string(userInput), "\r\n", "\n")
	text := strings.Split(input, "\n")
	
	var finalResult string

	for _, lines := range text {
		if lines == "" {
			finalResult += "\n"
			continue
		}

	finalResult += generateAscii(lines, realFont)
	}

	data := Pagedata{Title: "ASCII ART GENERATOR", Result: finalResult}
	renderTemplate(w, "index.html", data)
}


func generateAscii(text string, font []string) string{
	var result string
	
	for i :=1; i <= 8; i++{
		for _, char := range text {
			charIndex := (int(char)-32) * 9 + i
			if charIndex < len(font) {			
				result += font[charIndex]
			}
		}
		result += "\n"
	}
	return result
}



func renderTemplate(w http.ResponseWriter, Template string, data interface{}) {
	tmpl, err := template.ParseFiles(Template)
	if err != nil {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	tmpl.Execute(w, data)
}

