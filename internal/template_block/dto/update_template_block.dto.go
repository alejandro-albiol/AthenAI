package dto

// UpdateTemplateBlockDTO represents the data required to update an existing template block.
type UpdateTemplateBlockDTO struct {
	BlockName                *string `json:"block_name,omitempty"`
	BlockType                *string `json:"block_type,omitempty"`
	BlockOrder               *int    `json:"block_order,omitempty"`
	ExerciseCount            *int    `json:"exercise_count,omitempty"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
	Instructions             *string `json:"instructions,omitempty"`
	Reps                     *int    `json:"reps,omitempty"`
	Series                   *int    `json:"series,omitempty"`
	RestTimeSeconds          *int    `json:"rest_time_seconds,omitempty"`
}
