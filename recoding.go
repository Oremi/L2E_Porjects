package main

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Host string	`json: host`
	Port int	`json: port`
}

func parseJSON(data []byte) (string, int, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", 0, err
	}
	return cfg.Host, cfg.Port, nil
}

func main() {
	var data []byte = {port: 8080, host: "localhost"}
	host, port, err := (parseJSON(data))
	fmt.Println(host, port, err)
}



//func queryValidator(query string) string {
//	if r.URL.Query().Get(Query) != ""{
//		return "Query Not Found"
//	}
//}

func httpLogger(r, *http.Request){
	method := r.MethodPost
	query := r.URL.Query()
	cookie := r.Cookie("session.ID")

}

package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// TemplateCache holds our pre-parsed templates in memory RAM
type TemplateCache map[string]*template.Template

// NewTemplateCache scans a directory and compiles all templates at startup
func NewTemplateCache(dir string) (TemplateCache, error) {
	cache := make(TemplateCache)

	// 1. Find all top-level page templates (e.g., home.page.html)
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the clean filename (e.g., "home.page.html") to use as our map key
		name := filepath.Base(page)

		// 2. Parse the specific page template first
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
}

		// 3. Automatically append any shared global layout components (e.g., base.layout.html)
		// This merges your header, footer, or sidebar into the page template structure
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// 4. Store the fully compiled template bundle in the memory cache map
		cache[name] = ts
	}

	return cache, nil
}

// App holds our application dependencies, making the cache accessible to handlers
type App struct {
	templates TemplateCache
}

// homeHandler displays the homepage using our fast memory cache
func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the pre-parsed template directly from RAM
	ts, ok := app.templates["home.page.html"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// Dynamic data payload to pass to the HTML view
	data := map[string]interface{}{
		"Title": "Welcome to My Dashboard",
		"User":  "Alex",
	}

	// Render the template down to the network writer
	err := ts.Execute(w, data)
	if err != nil {
		log.Println("Execution error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	// Initialize and build the cache from your templates folder on startup
	// Assumes you have a folder structure containing your layout and page files
	cache, err := NewTemplateCache("./templates")
	if err != nil {
		log.Fatalf("Failed to initialize template cache: %v", err)
	}

	app := &App{templates: cache}

	http.HandleFunc("/", app.homeHandler)

	log.Println("Server booting up on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

