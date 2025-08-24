package dto

type UpdateCustomWorkoutTemplateDTO struct {
	ID                       string  `json:"id"`
	Name                     *string `json:"name,omitempty"`
	Description              *string `json:"description,omitempty"`
	DifficultyLevel          *string `json:"difficulty_level,omitempty"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
	TargetAudience           *string `json:"target_audience,omitempty"`
}
