package main

import (
	"github.com/alejandro-albiol/athenai/config"
	"log"
	"net/http"
	"os"
)

func main() {
	config.LoadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT undefined, using 8080")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, AthenAI!"))
	})

	log.Printf("Server is running on port: %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
