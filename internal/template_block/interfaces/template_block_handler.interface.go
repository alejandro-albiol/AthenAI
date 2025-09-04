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
	GetTemplateBlockByID(w http.ResponseWriter, r *http.Request)
	ListTemplateBlocksByTemplateID(w http.ResponseWriter, r *http.Request)
	UpdateTemplateBlock(w http.ResponseWriter, r *http.Request)
	DeleteTemplateBlock(w http.ResponseWriter, r *http.Request)
}
