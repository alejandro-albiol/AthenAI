package interfaces

import "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"

type ExerciseMuscularGroupRepository interface {
	CreateLink(link *dto.ExerciseMuscularGroup) (*string, error)
	DeleteLink(id string) error
	RemoveAllLinksForExercise(exerciseID string) error
	FindByID(id string) (*dto.ExerciseMuscularGroup, error)
	FindByExerciseID(exerciseID string) ([]*dto.ExerciseMuscularGroup, error)
	FindByMuscularGroupID(muscularGroupID string) ([]*dto.ExerciseMuscularGroup, error)
}
