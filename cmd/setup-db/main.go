package main

import (
	"log"

	"github.com/alejandro-albiol/athenai/config"
	"github.com/alejandro-albiol/athenai/internal/database"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database connection
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup database tables
	err = database.CreatePublicTables(db)
	if err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}

	log.Println("Database setup completed successfully!")
}
