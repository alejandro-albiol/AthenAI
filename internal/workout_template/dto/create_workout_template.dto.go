package dto

// CreateWorkoutTemplateDTO represents the data required to create a new workout template.
type CreateWorkoutTemplateDTO struct {
	Name                     string  `json:"name"`
	Description              string `json:"description,omitempty"`
	DifficultyLevel          string  `json:"difficulty_level"`
	EstimatedDurationMinutes int     `json:"estimated_duration_minutes,omitempty"`
	TargetAudience           string  `json:"target_audience,omitempty"`
	CreatedBy                string  `json:"created_by"` // ID of the user creating the template
}
