package dto

// UpdateTemplateBlockDTO represents the data required to update an existing template block.
type UpdateTemplateBlockDTO struct {
	Name                     *string  `json:"name,omitempty"`
	Type                     *string  `json:"type,omitempty"`
	Order                    *int     `json:"order,omitempty"`
	ExerciseCount            *int     `json:"exercise_count,omitempty"`
	EstimatedDurationMinutes *int     `json:"estimated_duration_minutes,omitempty"`
	Instructions             *string  `json:"instructions,omitempty"`
}