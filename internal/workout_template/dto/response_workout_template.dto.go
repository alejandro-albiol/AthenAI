package dto

import "time"

// ResponseWorkoutTemplateDTO represents the response structure for a workout template.
type ResponseWorkoutTemplateDTO struct {
	ID                       string    `json:"id"`
	Name                     string    `json:"name"`
	Description              string    `json:"description"`
	DifficultyLevel          string    `json:"difficulty_level"`
	EstimatedDurationMinutes int       `json:"estimated_duration_minutes"`
	TargetAudience           string    `json:"target_audience"`
	CreatedBy                string    `json:"created_by"`
	IsActive                 bool      `json:"is_active"`
	IsPublic                 bool      `json:"is_public"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}
