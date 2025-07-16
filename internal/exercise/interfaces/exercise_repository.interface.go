package interfaces

import "github.com/alejandro-albiol/athenai/internal/exercise/dto"

type ExerciseRepository interface {
	// Exercise management
	CreateExercise(exercise dto.ExerciseCreationDTO) (string, error)
	GetExerciseByID(id string) (dto.ExerciseResponseDTO, error)
	GetAllExercises() ([]dto.ExerciseResponseDTO, error)
	UpdateExercise(id string, exercise dto.ExerciseUpdateDTO) (dto.ExerciseResponseDTO, error)
	DeleteExercise(id string) error
	GetExercisesByMuscularGroup(muscularGroups []string) ([]dto.ExerciseResponseDTO, error)
	GetExercisesByEquipment(equipment []string) ([]dto.ExerciseResponseDTO, error)
}
