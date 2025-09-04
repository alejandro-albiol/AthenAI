package repository

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/workout_template/dto"
)

type WorkoutTemplateRepository struct {
	db *sql.DB
}

func NewWorkoutTemplateRepository(db *sql.DB) *WorkoutTemplateRepository {
	return &WorkoutTemplateRepository{db: db}
}

// Create a new workout template in the database
func (r *WorkoutTemplateRepository) CreateWorkoutTemplate(template *dto.CreateWorkoutTemplateDTO) (*string, error) {
	query := `INSERT INTO public.workout_template (name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id string
	err := r.db.QueryRow(query, template.Name, template.Description, template.DifficultyLevel,
		template.EstimatedDurationMinutes, template.TargetAudience, template.CreatedBy).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// Get a workout template by its ID from the database
func (r *WorkoutTemplateRepository) GetWorkoutTemplateByID(id string) (*dto.ResponseWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE id = $1`
	var template dto.ResponseWorkoutTemplateDTO
	err := r.db.QueryRow(query, id).Scan(&template.ID, &template.Name, &template.Description,
		&template.DifficultyLevel, &template.EstimatedDurationMinutes,
		&template.TargetAudience, &template.CreatedBy, &template.IsActive,
		&template.IsPublic, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout template by ID: %w", err)
	}
	return &template, nil
}

// Get a workout template by its name from the database
func (r *WorkoutTemplateRepository) GetWorkoutTemplateByName(name string) (*dto.ResponseWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE name = $1`
	var template dto.ResponseWorkoutTemplateDTO
	err := r.db.QueryRow(query, name).Scan(&template.ID, &template.Name, &template.Description,
		&template.DifficultyLevel, &template.EstimatedDurationMinutes,
		&template.TargetAudience, &template.CreatedBy, &template.IsActive,
		&template.IsPublic, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout template by name: %w", err)
	}
	return &template, nil
}

// Get workout templates by difficulty from the database
func (r *WorkoutTemplateRepository) GetWorkoutTemplatesByDifficulty(difficulty string) ([]*dto.ResponseWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE difficulty_level = $1`
	rows, err := r.db.Query(query, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout templates by difficulty: %w", err)
	}
	defer rows.Close()

	var templates []*dto.ResponseWorkoutTemplateDTO
	for rows.Next() {
		var template dto.ResponseWorkoutTemplateDTO
		err := rows.Scan(&template.ID, &template.Name, &template.Description,
			&template.DifficultyLevel, &template.EstimatedDurationMinutes,
			&template.TargetAudience, &template.CreatedBy, &template.IsActive,
			&template.IsPublic, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout template: %w", err)
		}
		templates = append(templates, &template)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed during iteration over workout templates: %w", err)
	}

	return templates, nil
}

// Get workout templates by target audience from the database
func (r *WorkoutTemplateRepository) GetWorkoutTemplatesByTargetAudience(targetAudience string) ([]*dto.ResponseWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE target_audience = $1`
	rows, err := r.db.Query(query, targetAudience)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout templates by target audience: %w", err)
	}
	defer rows.Close()

	var templates []*dto.ResponseWorkoutTemplateDTO
	for rows.Next() {
		var template dto.ResponseWorkoutTemplateDTO
		err := rows.Scan(&template.ID, &template.Name, &template.Description,
			&template.DifficultyLevel, &template.EstimatedDurationMinutes,
			&template.TargetAudience, &template.CreatedBy, &template.IsActive,
			&template.IsPublic, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout template: %w", err)
		}
		templates = append(templates, &template)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed during iteration over workout templates: %w", err)
	}

	return templates, nil
}

// Get all workout templates from the database
func (r *WorkoutTemplateRepository) GetAllWorkoutTemplates() ([]*dto.ResponseWorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all workout templates: %w", err)
	}
	defer rows.Close()

	var templates []*dto.ResponseWorkoutTemplateDTO
	for rows.Next() {
		var template dto.ResponseWorkoutTemplateDTO
		err := rows.Scan(&template.ID, &template.Name, &template.Description,
			&template.DifficultyLevel, &template.EstimatedDurationMinutes,
			&template.TargetAudience, &template.CreatedBy, &template.IsActive,
			&template.IsPublic, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout template: %w", err)
		}
		templates = append(templates, &template)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed during iteration over workout templates: %w", err)
	}

	return templates, nil
}

// Update a workout template in the database
func (r *WorkoutTemplateRepository) UpdateWorkoutTemplate(id string, template *dto.UpdateWorkoutTemplateDTO) (*dto.ResponseWorkoutTemplateDTO, error) {
	query := `UPDATE public.workout_template SET name = $1, description = $2, difficulty_level = $3, estimated_duration_minutes = $4, target_audience = $5, is_active = $6, is_public = $7, updated_at = NOW() WHERE id = $8 RETURNING id`
	_, err := r.db.Exec(query, template.Name, template.Description, template.DifficultyLevel,
		template.EstimatedDurationMinutes, template.TargetAudience,
		template.IsActive, template.IsPublic, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update workout template: %w", err)
	}

	updatedWorkoutTemplate, err := r.GetWorkoutTemplateByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated workout template: %w", err)
	}

	return updatedWorkoutTemplate, nil
}

// Delete a workout template from the database
func (r *WorkoutTemplateRepository) DeleteWorkoutTemplate(id string) error {
	query := `DELETE FROM public.workout_template WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workout template: %w", err)
	}
	return nil
}
