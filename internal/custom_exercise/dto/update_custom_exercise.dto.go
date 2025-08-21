package dto

// CustomExerciseUpdateDTO is used for updating a custom exercise
type CustomExerciseUpdateDTO struct {
	Name            *string   `json:"name"`
	Synonyms        *[]string `json:"synonyms"`
	DifficultyLevel *string   `json:"difficulty_level"`
	ExerciseType    *string   `json:"exercise_type"`
	Instructions    *string   `json:"instructions"`
	VideoURL        *string   `json:"video_url"`
	ImageURL        *string   `json:"image_url"`
	MuscularGroups  *[]string `json:"muscular_groups"`
}
