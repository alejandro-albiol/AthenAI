package interfaces

import "github.com/alejandro-albiol/athenai/internal/template_block/dto"
// TemplateBlockHandlerInterface defines the interface for the template block handler
// This should be used for dependency injection in the router layer
// to ensure strict module boundaries and testability.
type TemplateBlockHandlerInterface interface {
	// CreateTemplateBlock creates a new template block
	CreateTemplateBlock(block dto.CreateTemplateBlockDTO) error
	// GetTemplateBlockByID retrieves a template block by its ID
	GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error)
	// ListTemplateBlocksByTemplateID lists all template blocks for a given template ID
	ListTemplateBlocksByTemplateID(templateID string) ([]dto.TemplateBlockDTO, error)
	// UpdateTemplateBlock updates an existing template block
	UpdateTemplateBlock(id string, update dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error)
	// DeleteTemplateBlock deletes a template block by its ID
	DeleteTemplateBlock(id string) error
}