package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"

// CustomExerciseService defines business logic for custom exercises
//go:generate mockery --name=CustomExerciseService

type CustomExerciseService interface {
	CreateCustomExercise(exercise dto.CustomExerciseCreationDTO) (string, error)
	UpdateCustomExercise(id string, update dto.CustomExerciseUpdateDTO) error
	GetCustomExerciseByID(id string) (dto.CustomExerciseResponseDTO, error)
	ListCustomExercises() ([]dto.CustomExerciseResponseDTO, error)
	DeleteCustomExercise(id string) error
}
