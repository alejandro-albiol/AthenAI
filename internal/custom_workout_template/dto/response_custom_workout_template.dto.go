package dto

type ResponseCustomWorkoutTemplateDTO struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	Description              string `json:"description"`
	DifficultyLevel          string `json:"difficulty_level"`
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes"`
	TargetAudience           string `json:"target_audience"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
}
