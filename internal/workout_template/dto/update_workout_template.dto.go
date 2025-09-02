package dto

// UpdateWorkoutTemplateDTO is the data transfer object for updating a workout template.
type UpdateWorkoutTemplateDTO struct {
	Name                     *string  `json:"name"`
	Description              *string  `json:"description,omitempty"`
	DifficultyLevel          *string  `json:"difficulty_level"`
	EstimatedDurationMinutes *int     `json:"estimated_duration_minutes,omitempty"`
	TargetAudience           *string  `json:"target_audience,omitempty"`
	IsActive                 *bool    `json:"is_active"`
	IsPublic                 *bool    `json:"is_public"`
}
