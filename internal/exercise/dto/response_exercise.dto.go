package dto

import (
	"time"

	"github.com/lib/pq"
)

type ExerciseResponseDTO struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Synonyms         pq.StringArray `json:"synonyms"`
	MuscularGroups   pq.StringArray `json:"muscular_groups"`
	EquipmentNeeded  pq.StringArray `json:"equipment_needed"`
	DifficultyLevel  string         `json:"difficulty_level"`
	ExerciseType     string         `json:"exercise_type"`
	Instructions     string         `json:"instructions"`
	MuscularGroupIDs []string       `json:"muscular_group_ids"`
	EquipmentIDs     []string       `json:"equipment_ids"`
	VideoURL         *string        `json:"video_url"`
	ImageURL         *string        `json:"image_url"`
	CreatedBy        string         `json:"created_by"`
	IsActive         bool           `json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
