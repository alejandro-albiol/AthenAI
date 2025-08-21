package dto

import "github.com/alejandro-albiol/athenai/internal/exercise/enum"

type ExerciseCreationDTO struct {
	Name            string               `json:"name" validate:"required"`
	Synonyms        []string             `json:"synonyms"`
	MuscularGroups  []string             `json:"muscular_groups" validate:"required"`
	Equipment       []string             `json:"equipment"`
	DifficultyLevel enum.DifficultyLevel `json:"difficulty_level" validate:"required"`
	ExerciseType    enum.ExerciseType    `json:"exercise_type" validate:"required"`
	Instructions    string               `json:"instructions" validate:"required"`
	VideoURL        *string              `json:"video_url"`
	ImageURL        *string              `json:"image_url"`
	CreatedBy       string               `json:"created_by" validate:"required"`
}
