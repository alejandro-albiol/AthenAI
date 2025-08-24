package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"

type CustomTemplateBlockRepository interface {
	Create(gymID string, block *dto.CreateCustomTemplateBlockDTO) (string, error)
	GetByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error)
	ListByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error)
	Update(gymID string, block *dto.UpdateCustomTemplateBlockDTO) error
	Delete(gymID, id string) error
}
