package dto

type CreateCustomWorkoutInstanceDTO struct {
	Name                     string  `json:"name"`
	Description              string  `json:"description"`
	TemplateSource           string  `json:"template_source"`
	PublicTemplateID         *string `json:"public_template_id,omitempty"`
	GymTemplateID            *string `json:"gym_template_id,omitempty"`
	DifficultyLevel          string  `json:"difficulty_level"`
	EstimatedDurationMinutes int     `json:"estimated_duration_minutes"`
}
