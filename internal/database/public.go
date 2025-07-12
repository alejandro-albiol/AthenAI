package database

import (
	"database/sql"
	"fmt"
)

func CreatePublicTables(db *sql.DB) error {
	// Create gym table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.gym (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		domain TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL,
		address TEXT NOT NULL,
		phone TEXT NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		deleted_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		business_hours JSONB NOT NULL DEFAULT '[]',
		social_links JSONB NOT NULL DEFAULT '[]',
		payment_methods JSONB NOT NULL DEFAULT '[]'
	)`)
	if err != nil {
		return fmt.Errorf("failed to create gym table: %w", err)
	}
	fmt.Println("Gym table created successfully")

	// Create admin table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.admin (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create admin table: %w", err)
	}
	fmt.Println("Admin table created successfully")

	// Create exercise table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.exercise (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		synonyms TEXT[] NOT NULL DEFAULT '{}',
		muscular_groups TEXT[] NOT NULL DEFAULT '{}',
		equipment_needed TEXT[] NOT NULL DEFAULT '{}',
		difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
		exercise_type TEXT NOT NULL CHECK (exercise_type IN ('strength', 'cardio', 'flexibility', 'balance', 'functional')),
		instructions TEXT NOT NULL,
		video_url TEXT,
		image_url TEXT,
		ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
		ai_model_version TEXT,
		created_by UUID REFERENCES public.admin(id),
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create exercise table: %w", err)
	}
	fmt.Println("Exercise table created successfully")

	// Create muscular_group table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.muscular_group (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		body_part TEXT NOT NULL CHECK (body_part IN ('upper_body', 'lower_body', 'core', 'full_body')),
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create muscular_group table: %w", err)
	}
	fmt.Println("Muscular group table created successfully")

	// Create equipment table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.equipment (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		category TEXT NOT NULL CHECK (category IN ('free_weights', 'machines', 'cardio', 'accessories', 'bodyweight')),
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create equipment table: %w", err)
	}
	fmt.Println("Equipment table created successfully")

	return nil
}