package main

import (
	"embed"
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
)

//go:embed index.html
var content embed.FS

type Todo struct {
    UserID	int 	`json:"userId"`
    ID		int 	`json:"id"`
    Title	string	`json:"title"`
    Completed	bool	`json:"completed"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", todoHandler)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
    url := "https://jsonplaceholder.typicode.com/todos"
    todos, err := fetchTodos(url)
    if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }
	tmpl := template.Must(template.ParseFS(content, "index.html"))
	tmpl.Execute(w, todos)
}

func fetchTodos(url string) ([]*Todo, error) {
	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status: %s", resp.StatusCode)
	}

	// Read and Decode the body
	var data []*Todo
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
