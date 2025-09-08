package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"

type CustomMemberWorkoutService interface {
	CreateCustomMemberWorkout(gymID string, memberWorkout *dto.CreateCustomMemberWorkoutDTO) (*string, error)
	GetCustomMemberWorkoutByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error)
	ListCustomMemberWorkoutsByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error)
	UpdateCustomMemberWorkout(gymID string, memberWorkout *dto.UpdateCustomMemberWorkoutDTO) error
	DeleteCustomMemberWorkout(gymID, id string) error
}
