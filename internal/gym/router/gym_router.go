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
	r.Post("/", handler.CreateGym)                  // POST /gym
	r.Get("/", handler.GetAllGyms)                  // GET /gym
	r.Get("/{id}", handler.GetGymByID)              // GET /gym/{id}
	r.Get("/name/{name}", handler.GetGymByName)     // GET /gym/name/{name}
	r.Put("/{id}", handler.UpdateGym)               // PUT /gym/{id}
	r.Put("/{id}/activate", handler.SetGymActive)   // PUT /gym/{id}/activate
	r.Put("/{id}/deactivate", handler.SetGymActive) // PUT /gym/{id}/deactivate
	r.Delete("/{id}", handler.DeleteGym)            // DELETE /gym/{id}

	return r
}
