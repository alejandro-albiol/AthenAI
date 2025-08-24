package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"

type CustomWorkoutInstanceService interface {
	CreateCustomWorkoutInstance(gymID string, instance *dto.CreateCustomWorkoutInstanceDTO) (string, error)
	GetCustomWorkoutInstanceByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error)
	ListCustomWorkoutInstances(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	UpdateCustomWorkoutInstance(gymID string, instance *dto.UpdateCustomWorkoutInstanceDTO) error
	DeleteCustomWorkoutInstance(gymID, id string) error
}
