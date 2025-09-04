package interfaces

import "github.com/alejandro-albiol/athenai/internal/template_block/dto"

// TemplateBlockService defines the service interface for template blocks
// This should be used for dependency injection in handler and module layers
// to ensure strict module boundaries and testability.
type TemplateBlockService interface {
	// CreateTemplateBlock creates a new template block in the system.
	CreateTemplateBlock(block *dto.CreateTemplateBlockDTO) (*string, error)

	// GetTemplateBlockByID retrieves a template block by its unique ID.
	GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error)

	// ListTemplateBlocksByTemplateID lists all template blocks for a given template ID.
	ListTemplateBlocksByTemplateID(templateID string) ([]*dto.TemplateBlockDTO, error)

	// UpdateTemplateBlock updates an existing template block by its ID.
	UpdateTemplateBlock(id string, update *dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error)

	// DeleteTemplateBlock deletes a template block by its ID.
	DeleteTemplateBlock(id string) error
}
