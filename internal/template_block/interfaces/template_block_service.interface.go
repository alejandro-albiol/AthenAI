package interfaces

import "github.com/alejandro-albiol/athenai/internal/template_block/dto"

// TemplateBlockService defines the service interface for template blocks
// This should be used for dependency injection in handler and module layers
// to ensure strict module boundaries and testability.
type TemplateBlockService interface {
	CreateTemplateBlock(block dto.CreateTemplateBlockDTO) error
	GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error)
	ListTemplateBlocksByTemplateID(templateID string) ([]dto.TemplateBlockDTO, error)
	UpdateTemplateBlock(id string, update dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error)
	DeleteTemplateBlock(id string) error
}
