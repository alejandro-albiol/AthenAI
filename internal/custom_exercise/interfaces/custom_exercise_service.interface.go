package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"

// CustomExerciseService defines business logic for custom exercises
//go:generate mockery --name=CustomExerciseService

type CustomExerciseService interface {
	CreateCustomExercise(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error)
	UpdateCustomExercise(gymID string, id string, update *dto.CustomExerciseUpdateDTO) error
	GetCustomExerciseByID(gymID string, id string) (*dto.CustomExerciseResponseDTO, error)
	ListCustomExercises(gymID string) ([]*dto.CustomExerciseResponseDTO, error)
	DeleteCustomExercise(gymID string, id string) error
}
