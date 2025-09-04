package interfaces

import (
	"net/http"
)

// MuscularGroupHandlerInterface defines the interface for the muscular group handler
type MuscularGroupHandlerInterface interface {
	// CreateMuscularGroup creates a new muscular group
	CreateMuscularGroup(w http.ResponseWriter, r *http.Request)
	// GetMuscularGroup retrieves a muscular group by its ID
	GetMuscularGroup(w http.ResponseWriter, r *http.Request)
	// ListMuscularGroups lists all muscular groups
	ListMuscularGroups(w http.ResponseWriter, r *http.Request)
	// UpdateMuscularGroup updates an existing muscular group
	UpdateMuscularGroup(w http.ResponseWriter, r *http.Request)
	// DeleteMuscularGroup deletes a muscular group by its ID
	DeleteMuscularGroup(w http.ResponseWriter, r *http.Request)
}
