package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"

type CustomTemplateBlockService interface {
	CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (string, error)
	GetCustomTemplateBlockByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error)
	ListCustomTemplateBlocksByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error)
	UpdateCustomTemplateBlock(gymID string, block *dto.UpdateCustomTemplateBlockDTO) error
	DeleteCustomTemplateBlock(gymID, id string) error
}
