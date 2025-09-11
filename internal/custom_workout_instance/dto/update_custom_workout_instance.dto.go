package dto

type UpdateCustomWorkoutInstanceDTO struct {
	Name             *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description      *string `json:"description,omitempty" validate:"omitempty,max=500"`
	TemplateSource   *string `json:"template_source,omitempty" validate:"omitempty,oneof=public gym"`
	PublicTemplateID *string `json:"public_template_id,omitempty"`
	GymTemplateID    *string `json:"gym_template_id,omitempty"`
	// Note: DifficultyLevel and EstimatedDurationMinutes will be recalculated from exercises
}
