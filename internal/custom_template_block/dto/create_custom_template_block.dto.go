package dto

// CreateCustomTemplateBlockDTO is used for creating a custom template block
type CreateCustomTemplateBlockDTO struct {
	TemplateID               string `json:"template_id" validate:"required"`
	BlockName                string `json:"block_name" validate:"required"`
	BlockType                string `json:"block_type" validate:"required,oneof=warmup main core cardio cooldown custom"`
	BlockOrder               int    `json:"block_order" validate:"required,min=1"`
	ExerciseCount            int    `json:"exercise_count" validate:"required,min=1"`
	EstimatedDurationMinutes *int   `json:"estimated_duration_minutes,omitempty" validate:"omitempty,min=1"`
	Instructions             string `json:"instructions,omitempty"`
	Reps                     *int   `json:"reps,omitempty" validate:"omitempty,min=1"`
	Series                   *int   `json:"series,omitempty" validate:"omitempty,min=1"`
	RestTimeSeconds          *int   `json:"rest_time_seconds,omitempty" validate:"omitempty,min=0"`
	CreatedBy                string `json:"created_by" validate:"required"`
}
