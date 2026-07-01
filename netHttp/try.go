package main

import (
	"fmt"
	"net/http"
)

type PageData struct{
	Username string
	Age	int
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 Page Not Found", http.StatusNotFound)
			return
		}
		fmt.Println("Server: homeHandler")
		fmt.Fprintf(w, "This is a simple session with handlers")
	})
	fmt.Println("Server started at http://localhost:3000")
	fmt.Println(http.ListenAndServe(":3000", nil))
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w,r)
		return
	}
	fmt.Println("server: homeHandler session started...")

	fmt.Fprintf(w, "This is a sample request sent")
}
