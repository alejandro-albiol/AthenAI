package dto

type CreateCustomTemplateBlockDTO struct {
	TemplateID               string `json:"template_id"`
	BlockName                string `json:"block_name"`
	BlockType                string `json:"block_type"`
	BlockOrder               int    `json:"block_order"`
	ExerciseCount            int    `json:"exercise_count"`
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes"`
	Instructions             string `json:"instructions"`
}
