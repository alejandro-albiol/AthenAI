package dto

// UpdateCustomTemplateBlockDTO is used for updating a custom template block
type UpdateCustomTemplateBlockDTO struct {
	ID                       *string `json:"id"`
	BlockName                *string `json:"block_name,omitempty" validate:"omitempty"`
	BlockType                *string `json:"block_type,omitempty" validate:"omitempty,oneof=warmup main core cardio cooldown custom"`
	BlockOrder               *int    `json:"block_order,omitempty" validate:"omitempty,min=1"`
	ExerciseCount            *int    `json:"exercise_count,omitempty" validate:"omitempty,min=1"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty" validate:"omitempty,min=1"`
	Instructions             *string `json:"instructions,omitempty"`
	Reps                     *int    `json:"reps,omitempty" validate:"omitempty,min=1"`
	Series                   *int    `json:"series,omitempty" validate:"omitempty,min=1"`
	RestTimeSeconds          *int    `json:"rest_time_seconds,omitempty" validate:"omitempty,min=0"`
	IsActive                 *bool   `json:"is_active,omitempty"`
}
