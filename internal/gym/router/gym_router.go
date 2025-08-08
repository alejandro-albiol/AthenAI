package router

import (
	"net/http"

	gyminterfaces "github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewGymRouter(handler gyminterfaces.GymHandler) http.Handler {
	r := chi.NewRouter()

	// Auth middleware is applied globally at the API level
	// All routes here already have authenticated user context
	// Authorization logic is handled in the handlers themselves

	// Gym CRUD endpoints
	r.Post("/", handler.CreateGym)                  // POST /gyms
	r.Get("/", handler.GetAllGyms)                  // GET /gyms
	r.Get("/{id}", handler.GetGymByID)              // GET /gyms/{id}
	r.Get("/name/{name}", handler.GetGymByName)     // GET /gyms/name/{name}
	r.Put("/{id}", handler.UpdateGym)               // PUT /gyms/{id}
	r.Put("/{id}/activate", handler.SetGymActive)   // PUT /gyms/{id}/activate
	r.Put("/{id}/deactivate", handler.SetGymActive) // PUT /gyms/{id}/deactivate
	r.Delete("/{id}", handler.DeleteGym)            // DELETE /gyms/{id}

	return r
}
