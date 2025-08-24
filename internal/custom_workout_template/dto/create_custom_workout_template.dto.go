package dto

type CreateCustomWorkoutTemplateDTO struct {
	Name                     string `json:"name"`
	Description              string `json:"description"`
	DifficultyLevel          string `json:"difficulty_level"`
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes"`
	TargetAudience           string `json:"target_audience"`
}
