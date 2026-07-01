package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve all files in the current directory
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	fmt.Println("Server starting at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

