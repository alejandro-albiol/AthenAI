package dto

type CreateCustomWorkoutInstanceDTO struct {
	Name             string  `json:"name" validate:"required,min=3,max=100"`
	Description      string  `json:"description" validate:"max=500"`
	TemplateSource   string  `json:"template_source" validate:"required,oneof=public gym"`
	PublicTemplateID *string `json:"public_template_id,omitempty"`
	GymTemplateID    *string `json:"gym_template_id,omitempty"`
	// Note: DifficultyLevel and EstimatedDurationMinutes will be calculated from exercises
	// CreatedBy will be extracted from JWT token in the handler
}
