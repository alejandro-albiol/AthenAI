package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/interfaces"
)

type CustomMemberWorkoutService struct {
	repository interfaces.CustomMemberWorkoutRepository
}

func NewCustomMemberWorkoutService(repo interfaces.CustomMemberWorkoutRepository) *CustomMemberWorkoutService {
	return &CustomMemberWorkoutService{repository: repo}
}

func (s *CustomMemberWorkoutService) CreateCustomMemberWorkout(gymID string, memberWorkout *dto.CreateCustomMemberWorkoutDTO) (string, error) {
	return s.repository.Create(gymID, memberWorkout)
}

func (s *CustomMemberWorkoutService) GetCustomMemberWorkoutByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
	return s.repository.GetByID(gymID, id)
}

func (s *CustomMemberWorkoutService) ListCustomMemberWorkoutsByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error) {
	return s.repository.ListByMemberID(gymID, memberID)
}

func (s *CustomMemberWorkoutService) UpdateCustomMemberWorkout(gymID string, memberWorkout *dto.UpdateCustomMemberWorkoutDTO) error {
	return s.repository.Update(gymID, memberWorkout)
}

func (s *CustomMemberWorkoutService) DeleteCustomMemberWorkout(gymID, id string) error {
	return s.repository.Delete(gymID, id)
}
