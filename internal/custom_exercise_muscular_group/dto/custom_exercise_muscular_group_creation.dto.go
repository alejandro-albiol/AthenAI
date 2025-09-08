package dto

type CustomExerciseMuscularGroupCreationDTO struct {
	CustomExerciseID string `json:"custom_exercise_id" validate:"required"`
	MuscularGroupID  string `json:"muscular_group_id" validate:"required"`
}
