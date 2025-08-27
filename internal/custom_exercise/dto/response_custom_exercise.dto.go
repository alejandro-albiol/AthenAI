package dto

import "time"

// CustomExerciseResponseDTO is used for returning a custom exercise
type CustomExerciseResponseDTO struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Synonyms        []string   `json:"synonyms"`
	DifficultyLevel string     `json:"difficulty_level"`
	ExerciseType    string     `json:"exercise_type"`
	Instructions    string     `json:"instructions"`
	VideoURL        string     `json:"video_url"`
	ImageURL        string     `json:"image_url"`
	MuscularGroups  []string   `json:"muscular_groups"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	IsActive        bool       `json:"is_active"`
}
