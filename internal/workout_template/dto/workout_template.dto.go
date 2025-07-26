package dto

type WorkoutTemplateDTO struct {
	ID                       string  `json:"id"`
	Name                     string  `json:"name"`
	Description              *string `json:"description,omitempty"`
	DifficultyLevel          string  `json:"difficulty_level"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
	TargetAudience           *string `json:"target_audience,omitempty"`
	CreatedBy                string  `json:"created_by"`
	IsActive                 bool    `json:"is_active"`
	IsPublic                 bool    `json:"is_public"`
	CreatedAt                string  `json:"created_at"`
	UpdatedAt                string  `json:"updated_at"`
}
