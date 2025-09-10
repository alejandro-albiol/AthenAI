package dto

type TemplateBlockDTO struct {
	ID                       string  `json:"id"`
	TemplateID               string  `json:"template_id"`
	BlockName                string  `json:"block_name"`
	BlockType                string  `json:"block_type"`
	BlockOrder               int     `json:"block_order"`
	ExerciseCount            int     `json:"exercise_count"`
	EstimatedDurationMinutes *int    `json:"estimated_duration_minutes,omitempty"`
	Instructions             *string `json:"instructions,omitempty"`
	Reps                     *int    `json:"reps,omitempty"`
	Series                   *int    `json:"series,omitempty"`
	RestTimeSeconds          *int    `json:"rest_time_seconds,omitempty"`
	CreatedAt                string  `json:"created_at"`
	CreatedBy                string  `json:"created_by"`
}
