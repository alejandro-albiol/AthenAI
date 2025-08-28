package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"

type CustomExerciseMuscularGroupHandler interface {
	CreateLink(gymID string, link dto.CustomExerciseMuscularGroup) error
	DeleteLink(gymID, id string) error
	RemoveAllLinksForExercise(gymID, id string) error
	GetLinkByID(gymID, id string) (dto.CustomExerciseMuscularGroup, error)
	GetLinksByCustomExerciseID(gymID, customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error)
	GetLinksByMuscularGroupID(gymID, muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error)
}
