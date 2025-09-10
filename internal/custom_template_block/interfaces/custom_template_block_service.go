package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"

// CustomTemplateBlockService defines business logic for custom template blocks
//go:generate mockery --name=CustomTemplateBlockService

type CustomTemplateBlockService interface {
	CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (*string, error)
	UpdateCustomTemplateBlock(gymID string, id string, update *dto.UpdateCustomTemplateBlockDTO) (*dto.ResponseCustomTemplateBlockDTO, error)
	GetCustomTemplateBlockByID(gymID string, id string) (*dto.ResponseCustomTemplateBlockDTO, error)
	ListCustomTemplateBlocksByTemplateID(gymID string, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error)
	ListCustomTemplateBlocks(gymID string) ([]*dto.ResponseCustomTemplateBlockDTO, error)
	DeleteCustomTemplateBlock(gymID string, id string) error
}
