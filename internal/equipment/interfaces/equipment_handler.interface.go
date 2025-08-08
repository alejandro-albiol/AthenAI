package interfaces

import (
	"net/http"
)

// EquipmentHandlerInterface defines the interface for the equipment handler
// This should be used for dependency injection in the router layer
// to ensure strict module boundaries and testability.
type EquipmentHandlerInterface interface {
	// CreateEquipment creates a new equipment entry
	CreateEquipment(w http.ResponseWriter, r *http.Request)
	// GetEquipment retrieves an equipment entry by its ID
	GetEquipment(w http.ResponseWriter, r *http.Request)
	// ListEquipment lists all equipment entries
	ListEquipment(w http.ResponseWriter, r *http.Request)
	// UpdateEquipment updates an existing equipment entry
	UpdateEquipment(w http.ResponseWriter, r *http.Request)
	// DeleteEquipment deletes an equipment entry by its ID
	DeleteEquipment(w http.ResponseWriter, r *http.Request)
}
