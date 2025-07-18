package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"

type CustomExerciseMuscularGroupRepository interface {
	CreateLink(link dto.CustomExerciseMuscularGroup) (string, error)
	DeleteLink(id string) error
	FindByID(id string) (dto.CustomExerciseMuscularGroup, error)
	FindByCustomExerciseID(customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error)
	FindByMuscularGroupID(muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error)
}
