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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAllGyms(w, r)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateGym(w, r)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetGymByID(w, r)
	})

	r.Get("/name/{name}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetGymByName(w, r)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateGym(w, r)
	})

	r.Put("/{id}/activate", func(w http.ResponseWriter, r *http.Request) {
		handler.SetGymActive(w, r)
	})

	r.Put("/{id}/deactivate", func(w http.ResponseWriter, r *http.Request) {
		handler.SetGymActive(w, r)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteGym(w, r)
	})

	return r
}
