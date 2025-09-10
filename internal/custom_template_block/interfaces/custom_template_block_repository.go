package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"

// CustomTemplateBlockRepository defines DB operations for custom template blocks
// All methods must operate in the tenant schema

//go:generate mockery --name=CustomTemplateBlockRepository

type CustomTemplateBlockRepository interface {
	CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (*string, error)
	UpdateCustomTemplateBlock(gymID, id string, update *dto.UpdateCustomTemplateBlockDTO) error
	GetCustomTemplateBlockByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error)
	ListCustomTemplateBlocksByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error)
	ListCustomTemplateBlocks(gymID string) ([]*dto.ResponseCustomTemplateBlockDTO, error)
	DeleteCustomTemplateBlock(gymID, id string) error // soft delete
}
