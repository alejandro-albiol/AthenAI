package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alejandro-albiol/athenai/api"
	"github.com/alejandro-albiol/athenai/config"
	"github.com/alejandro-albiol/athenai/internal/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT undefined, using 8080")
	}

	// Initialize database connection
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Root router
	rootRouter := chi.NewRouter()

	// Setup Swagger at root level
	api.SetupSwagger(rootRouter)
	log.Println("Swagger setup at /swagger-ui/")

	// Mount API under /api/v1
	rootRouter.Mount("/api/v1", api.NewAPIModule(db))

	log.Printf("Server is running on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, rootRouter))
}
