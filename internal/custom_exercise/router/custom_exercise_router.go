package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/handler"
	"github.com/go-chi/chi/v5"
)

func NewCustomExerciseRouter(h *handler.CustomExerciseHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/custom-exercise", h.Create)
	r.Get("/custom-exercise/{id}", h.GetByID)
	r.Get("/custom-exercise", h.List)
	r.Put("/custom-exercise/{id}", h.Update)
	r.Delete("/custom-exercise/{id}", h.Delete)

	return r
}
