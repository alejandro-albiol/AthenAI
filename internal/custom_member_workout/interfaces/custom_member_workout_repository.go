package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"

type CustomMemberWorkoutRepository interface {
	Create(gymID string, memberWorkout *dto.CreateCustomMemberWorkoutDTO) (string, error)
	GetByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error)
	ListByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error)
	Update(gymID string, memberWorkout *dto.UpdateCustomMemberWorkoutDTO) error
	Delete(gymID, id string) error
}
