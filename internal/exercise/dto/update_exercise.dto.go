package dto

import "github.com/alejandro-albiol/athenai/internal/exercise/enum"

type ExerciseUpdateDTO struct {
	Name            *string              `json:"name"`
	Synonyms        []string             `json:"synonyms"`
	MuscularGroups  []string             `json:"muscular_groups"`
	Equipment       []string             `json:"equipment"`
	DifficultyLevel enum.DifficultyLevel `json:"difficulty_level" validate:"omitempty,oneof=beginner intermediate advanced"`
	ExerciseType    enum.ExerciseType    `json:"exercise_type" validate:"omitempty,oneof=strength cardio flexibility balance functional"`
	Instructions    *string              `json:"instructions"`
	VideoURL        *string              `json:"video_url"`
	ImageURL        *string              `json:"image_url"`
	IsActive        *bool                `json:"is_active"`
}
