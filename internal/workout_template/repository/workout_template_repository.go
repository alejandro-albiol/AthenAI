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
func (r *WorkoutTemplateRepository) Create(template dto.WorkoutTemplateDTO) error {
	query := `INSERT INTO public.workout_template (name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, template.Name, template.Description, template.DifficultyLevel,
		template.EstimatedDurationMinutes, template.TargetAudience, template.CreatedBy)
	return err
}

// Get a workout template by its ID from the database
func (r *WorkoutTemplateRepository) GetByID(id string) (*dto.WorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE id = $1`
	var template dto.WorkoutTemplateDTO
	err := r.db.QueryRow(query, id).Scan(&template.ID, &template.Name, &template.Description,
		&template.DifficultyLevel, &template.EstimatedDurationMinutes,
		&template.TargetAudience, &template.CreatedBy, &template.IsActive,
		&template.IsPublic, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout template by ID: %w", err)
	}

	return &template, nil
}

// Get all workout templates from the database
func (r *WorkoutTemplateRepository) GetAll() ([]dto.WorkoutTemplateDTO, error) {
	query := `SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all workout templates: %w", err)
	}
	defer rows.Close()

	var templates []dto.WorkoutTemplateDTO
	for rows.Next() {
		var template dto.WorkoutTemplateDTO
		err := rows.Scan(&template.ID, &template.Name, &template.Description,
			&template.DifficultyLevel, &template.EstimatedDurationMinutes,
			&template.TargetAudience, &template.CreatedBy, &template.IsActive,
			&template.IsPublic, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout template: %w", err)
		}
		templates = append(templates, template)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed during iteration over workout templates: %w", err)
	}

	return templates, nil
}

// Update a workout template in the database
func (r *WorkoutTemplateRepository) Update(template dto.WorkoutTemplateDTO) error {
	query := `UPDATE public.workout_template SET name = $1, description = $2, difficulty_level = $3, estimated_duration_minutes = $4, target_audience = $5, created_by = $6, is_active = $7, is_public = $8, updated_at = NOW() WHERE id = $9`
	_, err := r.db.Exec(query, template.Name, template.Description, template.DifficultyLevel,
		template.EstimatedDurationMinutes, template.TargetAudience, template.CreatedBy,
		template.IsActive, template.IsPublic, template.ID)
	if err != nil {
		return fmt.Errorf("failed to update workout template: %w", err)
	}

	return nil
}

// Delete a workout template from the database
func (r *WorkoutTemplateRepository) Delete(id string) error {
	query := `DELETE FROM public.workout_template WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workout template: %w", err)
	}

	return nil
}
