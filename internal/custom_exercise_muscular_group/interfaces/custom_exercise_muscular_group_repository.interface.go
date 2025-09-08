package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"

type CustomExerciseMuscularGroupRepository interface {
	CreateLink(gymID string, link *dto.CustomExerciseMuscularGroupCreationDTO) (*string, error)
	DeleteLink(gymID, id string) error
	RemoveAllLinksForExercise(gymID, customExerciseID string) error
	FindByID(gymID, id string) (*dto.CustomExerciseMuscularGroup, error)
	FindByCustomExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error)
	FindByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.CustomExerciseMuscularGroup, error)
}
