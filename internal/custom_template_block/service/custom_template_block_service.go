package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/interfaces"
)

type CustomTemplateBlockServiceImpl struct {
	Repo interfaces.CustomTemplateBlockRepository
}

func NewCustomTemplateBlockService(repo interfaces.CustomTemplateBlockRepository) *CustomTemplateBlockServiceImpl {
	return &CustomTemplateBlockServiceImpl{Repo: repo}
}

func (s *CustomTemplateBlockServiceImpl) CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (string, error) {
	return s.Repo.Create(gymID, block)
}

func (s *CustomTemplateBlockServiceImpl) GetCustomTemplateBlockByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error) {
	return s.Repo.GetByID(gymID, id)
}

func (s *CustomTemplateBlockServiceImpl) ListCustomTemplateBlocksByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	return s.Repo.ListByTemplateID(gymID, templateID)
}

func (s *CustomTemplateBlockServiceImpl) UpdateCustomTemplateBlock(gymID string, block *dto.UpdateCustomTemplateBlockDTO) error {
	return s.Repo.Update(gymID, block)
}

func (s *CustomTemplateBlockServiceImpl) DeleteCustomTemplateBlock(gymID, id string) error {
	return s.Repo.Delete(gymID, id)
}
