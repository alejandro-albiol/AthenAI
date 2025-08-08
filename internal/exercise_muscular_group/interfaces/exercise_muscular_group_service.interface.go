package interfaces

import "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"

// Service contract for join table

type ExerciseMuscularGroupService interface {
	CreateLink(link dto.ExerciseMuscularGroup) (string, error)
	DeleteLink(id string) error
	GetLinkByID(id string) (dto.ExerciseMuscularGroup, error)
	GetLinksByExerciseID(exerciseID string) ([]dto.ExerciseMuscularGroup, error)
	GetLinksByMuscularGroupID(muscularGroupID string) ([]dto.ExerciseMuscularGroup, error)
}
