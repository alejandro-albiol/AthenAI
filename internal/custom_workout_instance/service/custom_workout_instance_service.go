package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/interfaces"
)

type CustomWorkoutInstanceService struct {
	Repo interfaces.CustomWorkoutInstanceRepository
}

func NewCustomWorkoutInstanceService(repo interfaces.CustomWorkoutInstanceRepository) *CustomWorkoutInstanceService {
	return &CustomWorkoutInstanceService{Repo: repo}
}

func (s *CustomWorkoutInstanceService) CreateCustomWorkoutInstance(gymID string, instance *dto.CreateCustomWorkoutInstanceDTO) (string, error) {
	return s.Repo.Create(gymID, instance)
}

func (s *CustomWorkoutInstanceService) GetCustomWorkoutInstanceByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
	return s.Repo.GetByID(gymID, id)
}

func (s *CustomWorkoutInstanceService) ListCustomWorkoutInstances(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	return s.Repo.List(gymID)
}

func (s *CustomWorkoutInstanceService) UpdateCustomWorkoutInstance(gymID string, instance *dto.UpdateCustomWorkoutInstanceDTO) error {
	return s.Repo.Update(gymID, instance)
}

func (s *CustomWorkoutInstanceService) DeleteCustomWorkoutInstance(gymID, id string) error {
	return s.Repo.Delete(gymID, id)
}
