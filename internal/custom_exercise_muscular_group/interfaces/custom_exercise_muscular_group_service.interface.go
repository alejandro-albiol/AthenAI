package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"

type CustomExerciseMuscularGroupService interface {
	CreateLink(link dto.CustomExerciseMuscularGroup) error
	DeleteLink(id string) error
	GetLinkByID(id string) (dto.CustomExerciseMuscularGroup, error)
	GetLinksByCustomExerciseID(customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error)
	GetLinksByMuscularGroupID(muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error)
}
