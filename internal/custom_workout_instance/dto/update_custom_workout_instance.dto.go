package dto

type UpdateCustomWorkoutInstanceDTO struct {
	ID                       string  `json:"id"`
	Name                     *string `json:"name,omitempty"`
	Description              *string `json:"description,omitempty"`
	TemplateSource           *string `json:"template_source,omitempty"`
	PublicTemplateID         *string `json:"public_template_id,omitempty"`
	GymTemplateID            *string `json:"gym_template_id,omitempty"`
	DifficultyLevel          *string `json:"difficulty_level,omitempty"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
}
