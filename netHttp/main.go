package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)
	mux.HandleFunc("/hello", helloHandler)

	log.Println("Server is live at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	
	log.Println("server: Handler started")
	log.Println("This is in the Handler() function")
	fmt.Fprintf(w, "i dont know what this is")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Write data to the response
    fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}
