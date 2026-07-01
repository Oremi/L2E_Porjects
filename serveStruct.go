func main() {
	// 1. Initialize your Mux (router)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		// 2. Use an anonymous struct to parse specific JSON input
		var input struct {
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		// 3. Use a temporary struct to send a quick response
		json.NewEncoder(w).Encode(struct {
			Welcome string `json:"welcome"`
		}{
			Welcome: "Hello, " + input.Email,
		})
	})

	// 4. Configure the http.Server struct
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,            // Inject your mux here
		ReadTimeout:  5 * time.Second,  // Good for security
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}

