package interfaces

import (
	"net/http"
)

// TemplateBlockHandlerInterface defines the interface for the template block handler
// This should be used for dependency injection in the router layer
// to ensure strict module boundaries and testability.
type TemplateBlockHandlerInterface interface {
	// CreateTemplateBlock creates a new template block
	CreateTemplateBlock(w http.ResponseWriter, r *http.Request)
	// GetTemplateBlock retrieves a template block by its ID
	GetTemplateBlock(w http.ResponseWriter, r *http.Request)
	// ListTemplateBlocksByTemplateID lists all template blocks for a given template ID
	ListTemplateBlocksByTemplateID(w http.ResponseWriter, r *http.Request)
	// UpdateTemplateBlock updates an existing template block
	UpdateTemplateBlock(w http.ResponseWriter, r *http.Request)
	// DeleteTemplateBlock deletes a template block by its ID
	DeleteTemplateBlock(w http.ResponseWriter, r *http.Request)
}
