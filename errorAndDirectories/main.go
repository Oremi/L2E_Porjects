package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", RootHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is my content.")
	fmt.Fprintln(w, r.Header)
	
	defer r.Body.Close()
	
	body, err := io.ReadAll(r.Body)
	if err := nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, string(body))
}

func ResponseExampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintln(w, "Testing status code, Manually added a 200 Status OK.")
