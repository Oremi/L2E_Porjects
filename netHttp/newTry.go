package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HomeHandler)

	fmt.Println("Server started at http://localhost:8080")

	http.ListenAndServe(":8080", nil)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "This is a sample server session")

}

