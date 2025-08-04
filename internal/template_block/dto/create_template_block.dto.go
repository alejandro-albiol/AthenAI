package dto

// CreateTemplateBlockDTO represents the data required to create a new template block.
type CreateTemplateBlockDTO struct {
	TemplateID               string  `json:"template_id" validate:"required"`
	Name                     string  `json:"name" validate:"required"`
	Type                     string  `json:"type" validate:"required"`
	Order                    int     `json:"order" validate:"required"`
	ExerciseCount            int     `json:"exercise_count" validate:"required"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
	Instructions             *string `json:"instructions,omitempty"`
}
