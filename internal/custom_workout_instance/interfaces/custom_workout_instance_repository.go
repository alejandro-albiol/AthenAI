package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"

type CustomWorkoutInstanceRepository interface {
	Create(gymID string, instance dto.CreateCustomWorkoutInstanceDTO) (string, error)
	GetByID(gymID, id string) (dto.ResponseCustomWorkoutInstanceDTO, error)
	List(gymID string) ([]dto.ResponseCustomWorkoutInstanceDTO, error)
	Update(gymID string, instance dto.UpdateCustomWorkoutInstanceDTO) error
	Delete(gymID, id string) error
}
