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
		id := chi.URLParam(r, "id")
		handler.GetGymByID(w, r, id)
	})

	r.Get("/domain/{domain}", func(w http.ResponseWriter, r *http.Request) {
		domain := chi.URLParam(r, "domain")
		handler.GetGymByDomain(w, r, domain)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.UpdateGym(w, r, id)
	})

	r.Put("/{id}/activate", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.SetGymActive(w, r, id, true)
	})

	r.Put("/{id}/deactivate", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.SetGymActive(w, r, id, false)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.DeleteGym(w, r, id)
	})

	return r
}
