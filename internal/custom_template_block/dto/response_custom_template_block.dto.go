package dto

import "time"

// ResponseCustomTemplateBlockDTO is used for returning a custom template block
type ResponseCustomTemplateBlockDTO struct {
	ID                       string     `json:"id"`
	TemplateID               string     `json:"template_id"`
	BlockName                string     `json:"block_name"`
	BlockType                string     `json:"block_type"`
	BlockOrder               int        `json:"block_order"`
	ExerciseCount            int        `json:"exercise_count"`
	EstimatedDurationMinutes *int       `json:"estimated_duration_minutes,omitempty"`
	Instructions             string     `json:"instructions"`
	Reps                     *int       `json:"reps,omitempty"`
	Series                   *int       `json:"series,omitempty"`
	RestTimeSeconds          *int       `json:"rest_time_seconds,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty"`
	IsActive                 bool       `json:"is_active"`
	CreatedBy                string     `json:"created_by"`
}
