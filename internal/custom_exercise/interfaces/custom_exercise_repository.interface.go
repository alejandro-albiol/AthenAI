package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"

// CustomExerciseRepository defines DB operations for custom exercises
// All methods must operate in the tenant schema

//go:generate mockery --name=CustomExerciseRepository

type CustomExerciseRepository interface {
	CreateCustomExercise(exercise dto.CustomExerciseCreationDTO) (string, error)
	UpdateCustomExercise(id string, update dto.CustomExerciseUpdateDTO) error
	GetCustomExerciseByID(id string) (dto.CustomExerciseResponseDTO, error)
	ListCustomExercises() ([]dto.CustomExerciseResponseDTO, error)
	DeleteCustomExercise(id string) error // soft delete
}
