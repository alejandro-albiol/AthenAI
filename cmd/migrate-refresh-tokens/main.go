package main

import (
	"log"

	"github.com/alejandro-albiol/athenai/config"
	"github.com/alejandro-albiol/athenai/internal/database"
)

func main() {
	log.Println("=== AthenAI Database Migration: Refresh Tokens ===")
	log.Println("This will update the refresh_tokens table to use gym_id instead of gym_domain")

	// Load environment variables
	config.LoadEnv()

	// Initialize database connection
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Step 1: Dropping existing refresh_tokens table...")
	_, err = db.Exec("DROP TABLE IF EXISTS public.refresh_tokens")
	if err != nil {
		log.Fatalf("Failed to drop refresh_tokens table: %v", err)
	}
	log.Println("✅ Old refresh_tokens table dropped")

	log.Println("Step 2: Creating new refresh_tokens table with gym_id...")
	_, err = db.Exec(`
	CREATE TABLE public.refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		token TEXT NOT NULL UNIQUE,
		user_type VARCHAR(50) NOT NULL CHECK (user_type IN ('platform_admin', 'tenant_user')),
		gym_id UUID REFERENCES public.gym(id), -- NULL for platform admins, required for tenant users
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		
		-- Unique constraint: one refresh token per user per context
		UNIQUE(user_id, user_type, gym_id)
	)`)
	if err != nil {
		log.Fatalf("Failed to create new refresh_tokens table: %v", err)
	}
	log.Println("✅ New refresh_tokens table created")

	log.Println("Step 3: Creating indexes...")
	_, err = db.Exec(`
		CREATE INDEX idx_refresh_tokens_token ON public.refresh_tokens(token);
		CREATE INDEX idx_refresh_tokens_user ON public.refresh_tokens(user_id, user_type);
		CREATE INDEX idx_refresh_tokens_expires ON public.refresh_tokens(expires_at);
	`)
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	log.Println("✅ Indexes created")

	log.Println("🎉 Migration completed successfully!")
	log.Println("The refresh_tokens table now uses gym_id instead of gym_domain")
}
