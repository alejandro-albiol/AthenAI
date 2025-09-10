package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"

type CustomWorkoutExerciseRepository interface {
	Create(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error)
	GetByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error)
	ListByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error)
	ListByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error)
	ListByEquipmentID(gymID, equipmentID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error)
	Update(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error
	Delete(gymID, id string) error
}
