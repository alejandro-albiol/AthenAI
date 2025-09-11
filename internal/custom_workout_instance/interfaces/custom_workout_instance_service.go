package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"

type CustomWorkoutInstanceService interface {
	CreateCustomWorkoutInstance(gymID, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error)
	GetCustomWorkoutInstanceByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error)
	GetCustomWorkoutInstanceSummaryByID(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error)
	GetCustomWorkoutInstancesByUserID(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	GetCustomWorkoutInstanceSummariesByUserID(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error)
	GetLastCustomWorkoutInstancesByUserID(gymID, userID string, numberOfWorkouts int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	ListCustomWorkoutInstances(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	ListCustomWorkoutInstanceSummaries(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error)
	UpdateCustomWorkoutInstance(gymID, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error
	DeleteCustomWorkoutInstance(gymID, id string) error
}
