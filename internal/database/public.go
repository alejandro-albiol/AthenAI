package database

import (
	"database/sql"
	"fmt"
)

func CreatePublicTables(db *sql.DB) error {

	// 1. Admin table
	_, err := db.Exec(`
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

	// 2. Gym table
	_, err = db.Exec(`
	   CREATE TABLE IF NOT EXISTS public.gym (
		   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		   name TEXT NOT NULL,
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

	// 3. Muscular group table
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

	// 4. Equipment table
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

	// 5. Exercise table
	_, err = db.Exec(`
		  CREATE TABLE IF NOT EXISTS public.exercise (
				  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				  name TEXT NOT NULL,
				  synonyms TEXT[] NOT NULL,
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

	// 6. Join table for exercise <-> muscular_group
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

	// 7. Join table for exercise <-> equipment
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

	// 8. Workout template table
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
		   is_public BOOLEAN NOT NULL DEFAULT TRUE, -- If true, available to all gyms
		   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		   updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	   )`)
	if err != nil {
		return fmt.Errorf("failed to create workout_template table: %w", err)
	}
	fmt.Println("Workout template table created successfully")

	// 9. Template block table
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
		   reps INTEGER, -- Number of repetitions per exercise in this block
		   series INTEGER, -- Number of sets per exercise in this block
		   rest_time_seconds INTEGER, -- Rest time between sets in seconds
		   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		   created_by UUID REFERENCES public.admin(id),
       
		   -- Ensure unique order per template
		   UNIQUE(template_id, block_order)
	   )`)
	if err != nil {
		return fmt.Errorf("failed to create template_block table: %w", err)
	}
	fmt.Println("Template block table created successfully")

	// Add new fields to existing template_block table if they don't exist
	_, err = db.Exec(`
		ALTER TABLE public.template_block 
		ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES public.admin(id),
		ADD COLUMN IF NOT EXISTS reps INTEGER,
		ADD COLUMN IF NOT EXISTS series INTEGER,
		ADD COLUMN IF NOT EXISTS rest_time_seconds INTEGER
	`)
	if err != nil {
		return fmt.Errorf("failed to add new columns to template_block table: %w", err)
	}

	// 10. Refresh tokens table for JWT refresh token management
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

	// 11. Create indexes for refresh_token and template tables
	_, err = db.Exec(`
		   CREATE INDEX IF NOT EXISTS idx_refresh_token_token ON public.refresh_token(token);
		   CREATE INDEX IF NOT EXISTS idx_refresh_token_user ON public.refresh_token(user_id, user_type);
		   CREATE INDEX IF NOT EXISTS idx_refresh_token_expires ON public.refresh_token(expires_at);
		   CREATE INDEX IF NOT EXISTS idx_workout_template_active ON public.workout_template(is_active);
		   CREATE INDEX IF NOT EXISTS idx_workout_template_public ON public.workout_template(is_public);
		   CREATE INDEX IF NOT EXISTS idx_workout_template_difficulty ON public.workout_template(difficulty_level);
		   CREATE INDEX IF NOT EXISTS idx_template_block_template ON public.template_block(template_id);
		   CREATE INDEX IF NOT EXISTS idx_template_block_order ON public.template_block(template_id, block_order);
	   `)
	if err != nil {
		return fmt.Errorf("failed to create template indexes: %w", err)
	}
	fmt.Println("Template and refresh token indexes created successfully")

	return nil
}
