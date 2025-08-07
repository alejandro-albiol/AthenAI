package interfaces

import "net/http"

// GymHandler defines the interface for handling HTTP requests related to gyms.
// Each method corresponds to a specific API endpoint and handles the complete
// request-response cycle including error handling and response formatting.
type GymHandler interface {
	// CreateGym handles POST requests to create a new gym.
	// Reads gym data from request body and returns 201 on success.
	CreateGym(w http.ResponseWriter, r *http.Request)

	// GetGymByID handles GET requests to fetch a gym by its ID.
	// Returns 200 with gym data on success, 404 if not found.
	GetGymByID(w http.ResponseWriter, r *http.Request)

	// GetGymByName handles GET requests to fetch a gym by its name.
	// Returns 200 with gym data on success, 404 if not found.
	GetGymByName(w http.ResponseWriter, r *http.Request)

	// GetAllGyms handles GET requests to fetch all active gyms.
	// Returns 200 with array of gyms on success.
	GetAllGyms(w http.ResponseWriter, r *http.Request)

	// UpdateGym handles PUT/PATCH requests to update an existing gym.
	// Reads update data from request body and returns 200 on success.
	UpdateGym(w http.ResponseWriter, r *http.Request)

	// SetGymActive handles PATCH requests to change a gym's active status.
	// Returns 200 on success, 404 if gym not found.
	SetGymActive(w http.ResponseWriter, r *http.Request)

	// DeleteGym handles DELETE requests to remove a gym.
	// Returns 204 on success, 404 if gym not found.
	DeleteGym(w http.ResponseWriter, r *http.Request)
}
