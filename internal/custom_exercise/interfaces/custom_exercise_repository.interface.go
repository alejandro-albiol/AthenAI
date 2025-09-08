package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"

// CustomExerciseRepository defines DB operations for custom exercises
// All methods must operate in the tenant schema

//go:generate mockery --name=CustomExerciseRepository

type CustomExerciseRepository interface {
	CreateCustomExercise(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error)
	UpdateCustomExercise(gymID, id string, update *dto.CustomExerciseUpdateDTO) error
	GetCustomExerciseByID(gymID, id string) (*dto.CustomExerciseResponseDTO, error)
	ListCustomExercises(gymID string) ([]*dto.CustomExerciseResponseDTO, error)
	DeleteCustomExercise(gymID, id string) error // soft delete
}
