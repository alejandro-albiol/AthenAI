package router

import (
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/handler"
	"github.com/go-chi/chi/v5"
)

func NewCustomWorkoutTemplateRouter(h *handler.CustomWorkoutTemplateHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/custom-workout-template", h.Create)
	r.Get("/custom-workout-template/{id}", h.GetByID)
	r.Get("/custom-workout-template", h.List)
	r.Put("/custom-workout-template/{id}", h.Update)
	r.Delete("/custom-workout-template/{id}", h.Delete)
	return r
}
