package dto

type CustomExerciseMuscularGroup struct {
	ID               string `json:"id"`
	CustomExerciseID string `json:"custom_exercise_id"`
	MuscularGroupID  string `json:"muscular_group_id"`
}
