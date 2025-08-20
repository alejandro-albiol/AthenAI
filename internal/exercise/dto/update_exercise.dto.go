package dto

import "github.com/alejandro-albiol/athenai/internal/exercise/enum"

type ExerciseUpdateDTO struct {
	Name             *string  `json:"name"`
	Synonyms         []string `json:"synonyms"`
	MuscularGroups   []string `json:"muscular_groups"`
	EquipmentNeeded  []string `json:"equipment_needed"`
	DifficultyLevel  enum.DifficultyLevel  `json:"difficulty_level" validate:"omitempty,oneof=beginner intermediate advanced"`
	MuscularGroupIDs []string `json:"muscular_group_ids"`
	EquipmentIDs     []string `json:"equipment_ids"`
	ExerciseType     enum.ExerciseType  `json:"exercise_type" validate:"omitempty,oneof=strength cardio flexibility balance functional"`
	Instructions     *string  `json:"instructions"`
	VideoURL         *string  `json:"video_url"`
	ImageURL         *string  `json:"image_url"`
	IsActive         *bool    `json:"is_active"`
}
