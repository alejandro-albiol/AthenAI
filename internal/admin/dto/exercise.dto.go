package dto

import (
	"time"

	"github.com/lib/pq"
)

type ExerciseCreationDTO struct {
	Name            string   `json:"name" validate:"required"`
	Synonyms        []string `json:"synonyms"`
	MuscularGroups  []string `json:"muscular_groups" validate:"required"`
	EquipmentNeeded []string `json:"equipment_needed"`
	DifficultyLevel string   `json:"difficulty_level" validate:"required,oneof=beginner intermediate advanced"`
	ExerciseType    string   `json:"exercise_type" validate:"required,oneof=strength cardio flexibility balance functional"`
	Instructions    string   `json:"instructions" validate:"required"`
	VideoURL        *string  `json:"video_url"`
	ImageURL        *string  `json:"image_url"`
	CreatedBy       string   `json:"created_by" validate:"required"`
}

type ExerciseUpdateDTO struct {
	Name            *string  `json:"name"`
	Synonyms        []string `json:"synonyms"`
	MuscularGroups  []string `json:"muscular_groups"`
	EquipmentNeeded []string `json:"equipment_needed"`
	DifficultyLevel *string  `json:"difficulty_level" validate:"omitempty,oneof=beginner intermediate advanced"`
	ExerciseType    *string  `json:"exercise_type" validate:"omitempty,oneof=strength cardio flexibility balance functional"`
	Instructions    *string  `json:"instructions"`
	VideoURL        *string  `json:"video_url"`
	ImageURL        *string  `json:"image_url"`
	IsActive        *bool    `json:"is_active"`
}

type ExerciseResponseDTO struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Synonyms        pq.StringArray `json:"synonyms"`
	MuscularGroups  pq.StringArray `json:"muscular_groups"`
	EquipmentNeeded pq.StringArray `json:"equipment_needed"`
	DifficultyLevel string         `json:"difficulty_level"`
	ExerciseType    string         `json:"exercise_type"`
	Instructions    string         `json:"instructions"`
	VideoURL        *string        `json:"video_url"`
	ImageURL        *string        `json:"image_url"`
	CreatedBy       string         `json:"created_by"`
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
