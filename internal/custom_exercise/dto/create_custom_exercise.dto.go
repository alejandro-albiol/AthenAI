package dto

// CustomExerciseCreationDTO is used for creating a custom exercise
type CustomExerciseCreationDTO struct {
	Name            string   `json:"name" validate:"required"`
	Synonyms        []string `json:"synonyms"`
	DifficultyLevel string   `json:"difficulty_level"`
	ExerciseType    string   `json:"exercise_type"`
	Instructions    string   `json:"instructions"`
	VideoURL        string   `json:"video_url"`
	ImageURL        string   `json:"image_url"`
	MuscularGroups  []string `json:"muscular_groups"`
}
