package dto

type ResponseCustomWorkoutInstanceDTO struct {
	ID                       string  `json:"id"`
	Name                     string  `json:"name"`
	Description              string  `json:"description"`
	TemplateSource           string  `json:"template_source"`
	PublicTemplateID         *string `json:"public_template_id,omitempty"`
	GymTemplateID            *string `json:"gym_template_id,omitempty"`
	DifficultyLevel          string  `json:"difficulty_level"`
	EstimatedDurationMinutes int     `json:"estimated_duration_minutes"`
	CreatedAt                string  `json:"created_at"`
	UpdatedAt                string  `json:"updated_at"`
}
