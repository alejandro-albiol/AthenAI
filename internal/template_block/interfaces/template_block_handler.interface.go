package interfaces

import (
	"net/http"
)

// TemplateBlockHandlerInterface defines the interface for the template block handler
// This should be used for dependency injection in the router layer
// to ensure strict module boundaries and testability.
type TemplateBlockHandlerInterface interface {
	// CreateTemplateBlock creates a new template block
	CreateTemplateBlock(w http.ResponseWriter, r *http.Request) error
	// GetTemplateBlockByID retrieves a template block by its ID
	GetTemplateBlockByID(w http.ResponseWriter, r *http.Request) error
	// ListTemplateBlocksByTemplateID lists all template blocks for a given template ID
	ListTemplateBlocksByTemplateID(w http.ResponseWriter, r *http.Request) error
	// UpdateTemplateBlock updates an existing template block
	UpdateTemplateBlock(w http.ResponseWriter, r *http.Request) error
	// DeleteTemplateBlock deletes a template block by its ID
	DeleteTemplateBlock(w http.ResponseWriter, r *http.Request) error
}
