package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/interfaces"
)

type CustomWorkoutTemplateServiceImpl struct {
	Repo interfaces.CustomWorkoutTemplateRepository
}

func NewCustomWorkoutTemplateService(repo interfaces.CustomWorkoutTemplateRepository) *CustomWorkoutTemplateServiceImpl {
	return &CustomWorkoutTemplateServiceImpl{Repo: repo}
}

func (s *CustomWorkoutTemplateServiceImpl) CreateCustomWorkoutTemplate(gymID string, template *dto.CreateCustomWorkoutTemplateDTO) (string, error) {
	return s.Repo.Create(gymID, template)
}

func (s *CustomWorkoutTemplateServiceImpl) GetCustomWorkoutTemplateByID(gymID, id string) (*dto.ResponseCustomWorkoutTemplateDTO, error) {
	return s.Repo.GetByID(gymID, id)
}

func (s *CustomWorkoutTemplateServiceImpl) ListCustomWorkoutTemplates(gymID string) ([]*dto.ResponseCustomWorkoutTemplateDTO, error) {
	return s.Repo.List(gymID)
}

func (s *CustomWorkoutTemplateServiceImpl) UpdateCustomWorkoutTemplate(gymID string, template *dto.UpdateCustomWorkoutTemplateDTO) error {
	return s.Repo.Update(gymID, template)
}

func (s *CustomWorkoutTemplateServiceImpl) DeleteCustomWorkoutTemplate(gymID, id string) error {
	return s.Repo.Delete(gymID, id)
}
