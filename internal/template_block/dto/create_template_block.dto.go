package dto

// CreateTemplateBlockDTO represents the data required to create a new template block.
type CreateTemplateBlockDTO struct {
	TemplateID               string  `json:"template_id" validate:"required"`
	BlockName                string  `json:"block_name" validate:"required"`
	BlockType                string  `json:"block_type" validate:"required"`
	BlockOrder               int     `json:"block_order" validate:"required"`
	ExerciseCount            int     `json:"exercise_count" validate:"required"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
	Instructions             *string `json:"instructions,omitempty"`
}
