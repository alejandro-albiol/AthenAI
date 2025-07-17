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
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
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
		last_login_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create admin table: %w", err)
	}
	fmt.Println("Admin table created successfully")

	// Create exercise table (no arrays, use join tables for relations)
	_, err = db.Exec(`
	   CREATE TABLE IF NOT EXISTS public.exercise (
			   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			   name TEXT NOT NULL,
			   synonyms TEXT NOT NULL,
			   difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
			   exercise_type TEXT NOT NULL CHECK (exercise_type IN ('strength', 'cardio', 'flexibility', 'balance', 'functional')),
			   instructions TEXT NOT NULL,
			   video_url TEXT,
			   image_url TEXT,
			   created_by UUID REFERENCES public.admin(id),
			   is_active BOOLEAN NOT NULL DEFAULT TRUE,
			   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			   updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	   )`)
	if err != nil {
		return fmt.Errorf("failed to create exercise table: %w", err)
	}
	fmt.Println("Exercise table created successfully")

	// Create join table for exercise <-> muscular_group
	_, err = db.Exec(`
	   CREATE TABLE IF NOT EXISTS public.exercise_muscular_group (
			   exercise_id UUID NOT NULL REFERENCES public.exercise(id) ON DELETE CASCADE,
			   muscular_group_id UUID NOT NULL REFERENCES public.muscular_group(id) ON DELETE RESTRICT,
			   PRIMARY KEY (exercise_id, muscular_group_id)
	   )`)
	if err != nil {
		return fmt.Errorf("failed to create exercise_muscular_group table: %w", err)
	}
	fmt.Println("Exercise_muscular_group table created successfully")

	// Create join table for exercise <-> equipment
	_, err = db.Exec(`
	   CREATE TABLE IF NOT EXISTS public.exercise_equipment (
			   exercise_id UUID NOT NULL REFERENCES public.exercise(id) ON DELETE CASCADE,
			   equipment_id UUID NOT NULL REFERENCES public.equipment(id) ON DELETE RESTRICT,
			   PRIMARY KEY (exercise_id, equipment_id)
	   )`)
	if err != nil {
		return fmt.Errorf("failed to create exercise_equipment table: %w", err)
	}
	fmt.Println("Exercise_equipment table created successfully")

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

	// Create refresh_tokens table for JWT refresh token management
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.refresh_token (
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
		return fmt.Errorf("failed to create refresh_token table: %w", err)
	}
	fmt.Println("Refresh token table created successfully")

	// Create indexes for refresh_token table
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_token_token ON public.refresh_token(token);
		CREATE INDEX IF NOT EXISTS idx_refresh_token_user ON public.refresh_token(user_id, user_type);
		CREATE INDEX IF NOT EXISTS idx_refresh_token_expires ON public.refresh_token(expires_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create refresh_token indexes: %w", err)
	}
	fmt.Println("Refresh token indexes created successfully")

	// Create login_history table for audit trail (optional)
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.login_history (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		user_type VARCHAR(50) NOT NULL CHECK (user_type IN ('platform_admin', 'tenant_user')),
		gym_domain VARCHAR(255), -- NULL for platform admins
		ip_address INET,
		user_agent TEXT,
		login_successful BOOLEAN NOT NULL DEFAULT true,
		login_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create login_history table: %w", err)
	}
	fmt.Println("Login history table created successfully")

	// Create indexes for login_history table
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_login_history_user ON public.login_history(user_id, user_type);
		CREATE INDEX IF NOT EXISTS idx_login_history_time ON public.login_history(login_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create login_history indexes: %w", err)
	}
	fmt.Println("Login history indexes created successfully")

	// Create workout_template table for storing workout templates
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.workout_template (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		description TEXT,
		difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
		estimated_duration_minutes INTEGER,
		target_audience TEXT, -- e.g., 'weight_loss', 'muscle_building', 'endurance', 'general_fitness'
		created_by UUID REFERENCES public.admin(id),
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		is_public BOOLEAN NOT NULL DEFAULT FALSE, -- If true, available to all gyms
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create workout_template table: %w", err)
	}
	fmt.Println("Workout template table created successfully")

	// Create template_block table for organizing exercise slots within templates
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS public.template_block (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		template_id UUID NOT NULL REFERENCES public.workout_template(id) ON DELETE CASCADE,
		block_name TEXT NOT NULL, -- e.g., 'Pre-Warmup', 'Warmup', 'Main Block 1', 'Core', 'Cool Down'
		block_type TEXT NOT NULL CHECK (block_type IN ('warmup', 'main', 'core', 'cardio', 'cooldown', 'custom')),
		block_order INTEGER NOT NULL, -- Order of blocks in the template
		exercise_count INTEGER NOT NULL, -- Number of exercises for this block (e.g., 3 warmup exercises, 5 main exercises)
		estimated_duration_minutes INTEGER, -- Optional estimated time for this block
		instructions TEXT, -- Special instructions for this block type
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		
		-- Ensure unique order per template
		UNIQUE(template_id, block_order)
	)`)
	if err != nil {
		return fmt.Errorf("failed to create template_block table: %w", err)
	}
	fmt.Println("Template block table created successfully")

	// Create indexes for template tables
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_workout_template_active ON public.workout_template(is_active);
		CREATE INDEX IF NOT EXISTS idx_workout_template_public ON public.workout_template(is_public);
		CREATE INDEX IF NOT EXISTS idx_workout_template_difficulty ON public.workout_template(difficulty_level);
		CREATE INDEX IF NOT EXISTS idx_template_block_template ON public.template_block(template_id);
		CREATE INDEX IF NOT EXISTS idx_template_block_order ON public.template_block(template_id, block_order);
	`)
	if err != nil {
		return fmt.Errorf("failed to create template indexes: %w", err)
	}
	fmt.Println("Template indexes created successfully")

	return nil
}
