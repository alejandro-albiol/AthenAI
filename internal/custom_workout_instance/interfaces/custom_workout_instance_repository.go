package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"

type CustomWorkoutInstanceRepository interface {
	Create(gymID string, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error)
	GetByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error)
	GetSummaryByID(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error)
	GetByUserID(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	GetSummariesByUserID(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error)
	GetLastsByUserID(gymID, userID string, numberOfWorkouts int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	List(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	ListSummaries(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error)
	Update(gymID, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error
	Delete(gymID, id string) error
}
