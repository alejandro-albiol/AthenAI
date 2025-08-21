package dto

import (
	"time"
)

type ExerciseResponseDTO struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Synonyms        []string  `json:"synonyms"`
	MuscularGroups  []string  `json:"muscular_groups"`
	EquipmentNeeded []string  `json:"equipment_needed"`
	DifficultyLevel string    `json:"difficulty_level"`
	ExerciseType    string    `json:"exercise_type"`
	Instructions    string    `json:"instructions"`
	VideoURL        *string   `json:"video_url"`
	ImageURL        *string   `json:"image_url"`
	CreatedBy       string    `json:"created_by"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
