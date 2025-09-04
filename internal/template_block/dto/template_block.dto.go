package dto

type TemplateBlockDTO struct {
	ID                       string `json:"id"`
	TemplateID               string `json:"template_id"`
	Name                     string `json:"name"`
	Type                     string `json:"type"`
	Order                    int    `json:"order"`
	ExerciseCount            int    `json:"exercise_count"`
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes,omitempty"`
	Instructions             string `json:"instructions,omitempty"`
	CreatedAt                string `json:"created_at"`
}
