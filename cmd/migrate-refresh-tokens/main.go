package main

import (
	"log"

	"github.com/alejandro-albiol/athenai/config"
	"github.com/alejandro-albiol/athenai/internal/database"
)

func main() {
	log.Println("=== AthenAI Database Migration: Refresh Tokens ===")
	log.Println("This will update the refresh_token table to use gym_id instead of gym_domain")

	// Load environment variables
	config.LoadEnv()

	// Initialize database connection
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Step 1: Dropping existing refresh_token table...")
	_, err = db.Exec("DROP TABLE IF EXISTS public.refresh_token")
	if err != nil {
		log.Fatalf("Failed to drop refresh_token table: %v", err)
	}
	log.Println("✅ Old refresh_token table dropped")

	log.Println("Step 2: Creating new refresh_token table with gym_id...")
	_, err = db.Exec(`
	CREATE TABLE public.refresh_token (
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
		log.Fatalf("Failed to create new refresh_token table: %v", err)
	}
	log.Println("✅ New refresh_token table created")

	log.Println("Step 3: Creating indexes...")
	_, err = db.Exec(`
		CREATE INDEX idx_refresh_token_token ON public.refresh_token(token);
		CREATE INDEX idx_refresh_token_user ON public.refresh_token(user_id, user_type);
		CREATE INDEX idx_refresh_token_expires ON public.refresh_token(expires_at);
	`)
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	log.Println("✅ Indexes created")

	log.Println("🎉 Migration completed successfully!")
	log.Println("The refresh_token table now uses gym_id instead of gym_domain")
}
