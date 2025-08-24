package dto

type UpdateCustomTemplateBlockDTO struct {
	ID                       string  `json:"id"`
	BlockName                *string `json:"block_name,omitempty"`
	BlockType                *string `json:"block_type,omitempty"`
	BlockOrder               *int    `json:"block_order,omitempty"`
	ExerciseCount            *int    `json:"exercise_count,omitempty"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
	Instructions             *string `json:"instructions,omitempty"`
}
