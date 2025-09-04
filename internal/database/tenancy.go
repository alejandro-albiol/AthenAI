package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// CreateTenantSchema creates a new schema and users table for a tenant (gym)
func CreateTenantSchema(db *sql.DB, schemaName *string) error {
	schema := pq.QuoteIdentifier(*schemaName)
	// Helper to quote table names with schema
	qt := func(table string) string {
		return fmt.Sprintf("%s.%s", schema, pq.QuoteIdentifier(table))
	}
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

	// Create custom_exercise table for gym-specific exercises
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

	// Create join table for custom_exercise and muscular_group
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

	// Create join table for custom_exercise and equipment
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

	// Create custom_equipment table for gym-specific equipment
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

	// Create custom_workout_template table for gym-specific templates
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.custom_workout_template (
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
		return fmt.Errorf("failed to create custom_workout_template table: %w", err)
	}

	// Create custom_template_block table for gym-specific templates
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.custom_template_block (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			template_id UUID NOT NULL REFERENCES %s.custom_workout_template(id) ON DELETE CASCADE,
			block_name TEXT NOT NULL,
			block_type TEXT NOT NULL CHECK (block_type IN ('warmup', 'main', 'core', 'cardio', 'cooldown', 'custom')),
			block_order INTEGER NOT NULL,
			exercise_count INTEGER NOT NULL,
			estimated_duration_minutes INTEGER,
			instructions TEXT
		)
	`, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_template_block table: %w", err)
	}

	// Create custom_workout_instance table for actual workouts created from templates
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.custom_workout_instance (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			template_source TEXT NOT NULL CHECK (template_source IN ('public', 'gym')),
			public_template_id UUID,
			gym_template_id UUID,
			difficulty_level TEXT NOT NULL CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
			estimated_duration_minutes INTEGER
		)
	`, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_workout_instance table: %w", err)
	}

	// Create custom_workout_exercise table for exercises assigned to workout instances
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.custom_workout_exercise (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			workout_instance_id UUID NOT NULL REFERENCES %s.custom_workout_instance(id) ON DELETE CASCADE,
			exercise_source TEXT NOT NULL CHECK (exercise_source IN ('public', 'gym')),
			public_exercise_id UUID,
			gym_exercise_id UUID,
			block_name TEXT NOT NULL,
			exercise_order INTEGER NOT NULL,
			sets INTEGER,
			reps_min INTEGER,
			reps_max INTEGER,
			weight_kg DECIMAL(5,2),
			duration_seconds INTEGER,
			rest_seconds INTEGER,
			notes TEXT,
			CHECK (
				(exercise_source = 'public' AND public_exercise_id IS NOT NULL AND gym_exercise_id IS NULL) OR
				(exercise_source = 'gym' AND gym_exercise_id IS NOT NULL AND public_exercise_id IS NULL)
			),
			UNIQUE(workout_instance_id, block_name, exercise_order)
		)
	`, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_workout_exercise table: %w", err)
	}

	// Create custom_member_workout table for member workout sessions
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.custom_member_workout (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL,
			member_id UUID NOT NULL,
			workout_instance_id UUID NOT NULL REFERENCES %s.custom_workout_instance(id),
            
			scheduled_date DATE,
			started_at TIMESTAMP WITH TIME ZONE,
			completed_at TIMESTAMP WITH TIME ZONE,
			status TEXT NOT NULL CHECK (status IN ('scheduled', 'in_progress', 'completed', 'skipped', 'cancelled')),
            
			notes TEXT,
			rating INTEGER CHECK (rating >= 1 AND rating <= 5)
		)
	`, schema, schema))
	if err != nil {
		return fmt.Errorf("failed to create custom_member_workout table: %w", err)
	}

	// Create indexes for better performance
	quoteIdx := func(name string) string { return pq.QuoteIdentifier(name) }
	indexStmts := []string{
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(is_active);", quoteIdx("idx_"+*schemaName+"_custom_exercise_active"), qt("custom_exercise")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(exercise_type);", quoteIdx("idx_"+*schemaName+"_custom_exercise_type"), qt("custom_exercise")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(difficulty_level);", quoteIdx("idx_"+*schemaName+"_custom_exercise_difficulty"), qt("custom_exercise")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(is_active);", quoteIdx("idx_"+*schemaName+"_custom_equipment_active"), qt("custom_equipment")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(category);", quoteIdx("idx_"+*schemaName+"_custom_equipment_category"), qt("custom_equipment")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(is_active);", quoteIdx("idx_"+*schemaName+"_workout_template_active"), qt("workout_template")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(difficulty_level);", quoteIdx("idx_"+*schemaName+"_workout_template_difficulty"), qt("workout_template")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(template_id);", quoteIdx("idx_"+*schemaName+"_template_block_template"), qt("template_block")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(template_id, block_order);", quoteIdx("idx_"+*schemaName+"_template_block_order"), qt("template_block")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(public_template_id);", quoteIdx("idx_"+*schemaName+"_workout_instance_public_template"), qt("workout_instance")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(gym_template_id);", quoteIdx("idx_"+*schemaName+"_workout_instance_gym_template"), qt("workout_instance")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(workout_instance_id);", quoteIdx("idx_"+*schemaName+"_workout_exercise_instance"), qt("workout_exercise")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(public_exercise_id);", quoteIdx("idx_"+*schemaName+"_workout_exercise_public"), qt("workout_exercise")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(gym_exercise_id);", quoteIdx("idx_"+*schemaName+"_workout_exercise_gym"), qt("workout_exercise")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(workout_instance_id, block_name);", quoteIdx("idx_"+*schemaName+"_workout_exercise_block"), qt("workout_exercise")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(member_id);", quoteIdx("idx_"+*schemaName+"_member_workout_member"), qt("member_workout")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(workout_instance_id);", quoteIdx("idx_"+*schemaName+"_member_workout_instance"), qt("member_workout")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(status);", quoteIdx("idx_"+*schemaName+"_member_workout_status"), qt("member_workout")),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(scheduled_date);", quoteIdx("idx_"+*schemaName+"_member_workout_date"), qt("member_workout")),
	}
	for _, stmt := range indexStmts {
		fmt.Printf("[DEBUG] Executing index SQL: %s\n", stmt)
		if _, err := db.Exec(stmt); err != nil {
			fmt.Printf("[ERROR] Index creation failed: %v\n", err)
			return fmt.Errorf("failed to create tenant index: %w", err)
		}
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
