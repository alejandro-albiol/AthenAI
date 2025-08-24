package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"

type CustomWorkoutExerciseService interface {
	CreateCustomWorkoutExercise(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (string, error)
	GetCustomWorkoutExerciseByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error)
	ListCustomWorkoutExercisesByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error)
	UpdateCustomWorkoutExercise(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error
	DeleteCustomWorkoutExercise(gymID, id string) error
}
