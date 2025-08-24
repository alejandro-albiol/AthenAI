package dto

type ResponseCustomTemplateBlockDTO struct {
	ID                       string `json:"id"`
	TemplateID               string `json:"template_id"`
	BlockName                string `json:"block_name"`
	BlockType                string `json:"block_type"`
	BlockOrder               int    `json:"block_order"`
	ExerciseCount            int    `json:"exercise_count"`
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes"`
	Instructions             string `json:"instructions"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
}
