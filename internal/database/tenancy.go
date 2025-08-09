package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// CreateTenantSchema creates a new schema and users table for a tenant (gym)
func CreateTenantSchema(db *sql.DB, schemaName string) error {
	schema := pq.QuoteIdentifier(schemaName)
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Example: create a user table in the new schema
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.user (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			is_verified BOOLEAN NOT NULL DEFAULT false,
			is_active BOOLEAN NOT NULL DEFAULT true,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			description TEXT,
			training_phase TEXT CHECK (training_phase IN ('weight_loss', 'muscle_gain', 'cardio_improve', 'maintenance')),
			motivation TEXT CHECK (motivation IN ('medical_recommendation', 'self_improvement', 'competition', 'rehabilitation', 'wellbeing')),
			special_situation TEXT CHECK (special_situation IN ('pregnancy', 'post_partum', 'injury_recovery', 'chronic_condition', 'elderly_population', 'physical_limitation', 'none')),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`, schema))
	if err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	// Add more table creation as needed

	// Create custom_exercise table - gym-specific exercises that extend public library (no arrays, use join tables)
	_, err = db.Exec(fmt.Sprintf(`
			   CREATE TABLE IF NOT EXISTS %s.custom_exercise (
					   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					   created_by UUID NOT NULL,
					   name TEXT NOT NULL,
					   synonyms TEXT NOT NULL,
					   difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
					   exercise_type TEXT NOT NULL CHECK (exercise_type IN ('strength', 'cardio', 'flexibility', 'balance', 'functional')),
					   instructions TEXT NOT NULL,
					   video_url TEXT,
					   image_url TEXT,
					   is_active BOOLEAN NOT NULL DEFAULT TRUE
			   )
	   `, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_exercise table: %w", err)
	}

	// Create join table for custom_exercise <-> muscular_group
	_, err = db.Exec(fmt.Sprintf(`
			   CREATE TABLE IF NOT EXISTS %s.custom_exercise_muscular_group (
					   custom_exercise_id UUID NOT NULL REFERENCES %s.custom_exercise(id) ON DELETE CASCADE,
					   muscular_group_id UUID NOT NULL REFERENCES public.muscular_group(id) ON DELETE RESTRICT,
					   PRIMARY KEY (custom_exercise_id, muscular_group_id)
			   )
	   `, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_exercise_muscular_group table: %w", err)
	}

	// Create join table for custom_exercise <-> equipment
	_, err = db.Exec(fmt.Sprintf(`
			   CREATE TABLE IF NOT EXISTS %s.custom_exercise_equipment (
					   custom_exercise_id UUID NOT NULL REFERENCES %s.custom_exercise(id) ON DELETE CASCADE,
					   equipment_id UUID NOT NULL REFERENCES public.equipment(id) ON DELETE RESTRICT,
					   PRIMARY KEY (custom_exercise_id, equipment_id)
			   )
	   `, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_exercise_equipment table: %w", err)
	}

	// Create custom_equipment table - gym-specific equipment that extends public library
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.custom_equipment (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			category TEXT NOT NULL CHECK (category IN ('free_weights', 'machines', 'cardio', 'accessories', 'bodyweight', 'custom')),
			is_active BOOLEAN NOT NULL DEFAULT TRUE
		)
	`, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_equipment table: %w", err)
	}

	// Create workout_template table - gym-specific templates (similar to public.workout_template)
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.workout_template (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
			estimated_duration_minutes INTEGER,
			target_audience TEXT NOT NULL CHECK (target_audience IN ('weight_loss', 'muscle_building', 'endurance', 'strength', 'flexibility', 'general_fitness', 'rehabilitation')),
			is_active BOOLEAN NOT NULL DEFAULT TRUE
		)
	`, schema))
	if err != nil {
		return fmt.Errorf("failed to create workout_template table: %w", err)
	}

	// Create template_block table - blocks for gym-specific templates
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.template_block (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			template_id UUID NOT NULL REFERENCES %s.workout_template(id) ON DELETE CASCADE,
			block_name TEXT NOT NULL, -- e.g., 'Warmup', 'Main Block 1', 'Core', 'Cool Down'
			block_type TEXT NOT NULL CHECK (block_type IN ('warmup', 'main', 'core', 'cardio', 'cooldown', 'custom')),
			block_order INTEGER NOT NULL, -- Order of blocks in the template
			exercise_count INTEGER NOT NULL, -- Number of exercises for this block
			estimated_duration_minutes INTEGER, -- Optional estimated time for this block
			instructions TEXT -- Special instructions for this block type
		)
	`, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create template_block table: %w", err)
	}

	// Create workout_instance table - actual workouts created from templates
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.workout_instance (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			
			template_source TEXT NOT NULL CHECK (template_source IN ('public', 'gym')),
			public_template_id UUID, -- References public.workout_template if template_source='public'
			gym_template_id UUID,    -- References {gym}.workout_template if template_source='gym'
			
			-- Workout metadata
			difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
			estimated_duration_minutes INTEGER
		)
	`, schema))
	if err != nil {
		return fmt.Errorf("failed to create workout_instance table: %w", err)
	}

	// Create workout_exercise table - exercises assigned to workout instances
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.workout_exercise (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			workout_instance_id UUID NOT NULL REFERENCES %s.workout_instance(id) ON DELETE CASCADE,
			
			-- Exercise source (can be public or gym-specific)
			exercise_source TEXT NOT NULL CHECK (exercise_source IN ('public', 'gym')),
			public_exercise_id UUID,  -- References public.exercise if exercise_source='public'
			gym_exercise_id UUID,     -- References {gym}.custom_exercise if exercise_source='gym'
			
			-- Exercise configuration
			block_name TEXT NOT NULL, -- Which block this exercise belongs to
			exercise_order INTEGER NOT NULL, -- Order within the block
			
			-- Exercise parameters (can override defaults)
			sets INTEGER,
			reps_min INTEGER,
			reps_max INTEGER,
			weight_kg DECIMAL(5,2),
			duration_seconds INTEGER,
			rest_seconds INTEGER,
			notes TEXT,
			
			-- Ensure only one exercise source is set
			CHECK (
				(exercise_source = 'public' AND public_exercise_id IS NOT NULL AND gym_exercise_id IS NULL) OR
				(exercise_source = 'gym' AND gym_exercise_id IS NOT NULL AND public_exercise_id IS NULL)
			),
			
			-- Ensure unique order per block per workout
			UNIQUE(workout_instance_id, block_name, exercise_order)
		)
	`, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create workout_exercise table: %w", err)
	}

	// Create member_workout table - member workout sessions
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.member_workout (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			member_id UUID NOT NULL, -- References {gym}.user (the member)
			workout_instance_id UUID NOT NULL REFERENCES %s.workout_instance(id),
			
			-- Session details
			scheduled_date DATE,
			started_at TIMESTAMP WITH TIME ZONE,
			completed_at TIMESTAMP WITH TIME ZONE,
			status TEXT NOT NULL CHECK (status IN ('scheduled', 'in_progress', 'completed', 'skipped', 'cancelled')),
			
			-- Session metadata
			notes TEXT,
			rating INTEGER CHECK (rating >= 1 AND rating <= 5) -- Member's workout rating
		)
	`, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create member_workout table: %w", err)
	}

	// Create indexes for better performance
	_, err = db.Exec(fmt.Sprintf(`
		-- Indexes for custom_exercises
		CREATE INDEX IF NOT EXISTS idx_%s_custom_exercises_active ON %s.custom_exercises(is_active);
		CREATE INDEX IF NOT EXISTS idx_%s_custom_exercises_type ON %s.custom_exercises(exercise_type);
		CREATE INDEX IF NOT EXISTS idx_%s_custom_exercises_difficulty ON %s.custom_exercises(difficulty_level);
		
		-- Indexes for custom_equipment
		CREATE INDEX IF NOT EXISTS idx_%s_custom_equipment_active ON %s.custom_equipment(is_active);
		CREATE INDEX IF NOT EXISTS idx_%s_custom_equipment_category ON %s.custom_equipment(category);
		
		-- Indexes for workout_template
		CREATE INDEX IF NOT EXISTS idx_%s_workout_template_active ON %s.workout_template(is_active);
		CREATE INDEX IF NOT EXISTS idx_%s_workout_template_difficulty ON %s.workout_template(difficulty_level);
		
		-- Indexes for template_block
		CREATE INDEX IF NOT EXISTS idx_%s_template_block_template ON %s.template_block(template_id);
		CREATE INDEX IF NOT EXISTS idx_%s_template_block_order ON %s.template_block(template_id, block_order);
		
		-- Indexes for workout_instance
		CREATE INDEX IF NOT EXISTS idx_%s_workout_instance_active ON %s.workout_instance(is_active);
		CREATE INDEX IF NOT EXISTS idx_%s_workout_instance_public_template ON %s.workout_instance(public_template_id);
		CREATE INDEX IF NOT EXISTS idx_%s_workout_instance_gym_template ON %s.workout_instance(gym_template_id);
		
		-- Indexes for workout_exercise
		CREATE INDEX IF NOT EXISTS idx_%s_workout_exercise_instance ON %s.workout_exercise(workout_instance_id);
		CREATE INDEX IF NOT EXISTS idx_%s_workout_exercise_public ON %s.workout_exercise(public_exercise_id);
		CREATE INDEX IF NOT EXISTS idx_%s_workout_exercise_gym ON %s.workout_exercise(gym_exercise_id);
		CREATE INDEX IF NOT EXISTS idx_%s_workout_exercise_block ON %s.workout_exercise(workout_instance_id, block_name);
		
		-- Indexes for member_workout
		CREATE INDEX IF NOT EXISTS idx_%s_member_workout_member ON %s.member_workout(member_id);
		CREATE INDEX IF NOT EXISTS idx_%s_member_workout_instance ON %s.member_workout(workout_instance_id);
		CREATE INDEX IF NOT EXISTS idx_%s_member_workout_status ON %s.member_workout(status);
		CREATE INDEX IF NOT EXISTS idx_%s_member_workout_date ON %s.member_workout(scheduled_date);
	`,
		schema, schema, // custom_exercises
		schema, schema, schema, schema, // custom_exercises type & difficulty
		schema, schema, schema, schema, // custom_equipment
		schema, schema, schema, schema, // workout_template
		schema, schema, schema, schema, // template_block
		schema, schema, schema, schema, schema, schema, // workout_instance
		schema, schema, schema, schema, schema, schema, schema, schema, // workout_exercise
		schema, schema, schema, schema, schema, schema, schema, schema)) // member_workout
	if err != nil {
		return fmt.Errorf("failed to create tenant indexes: %w", err)
	}

	return nil
}

// DeleteTenantSchema completely removes a tenant's schema and all its data
// WARNING: This operation is irreversible and will delete ALL tenant data
func DeleteTenantSchema(db *sql.DB, schemaName string) error {
	schema := pq.QuoteIdentifier(schemaName)

	// Drop the entire schema with CASCADE to remove all objects within it
	_, err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
	if err != nil {
		return fmt.Errorf("failed to drop schema %s: %w", schemaName, err)
	}

	return nil
}

// ListTenantSchemas returns all tenant schemas in the database
func ListTenantSchemas(db *sql.DB) ([]string, error) {
	query := `
		SELECT schema_name 
		FROM information_schema.schemata 
		WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast', 'public')
		AND schema_name NOT LIKE 'pg_%'
		ORDER BY schema_name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas: %w", err)
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}

	return schemas, nil
}
