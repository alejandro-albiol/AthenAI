package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewCustomWorkoutInstanceRouter(handler interfaces.CustomWorkoutInstanceHandler) http.Handler {
	router := chi.NewRouter()
	router.Post("/custom-workout-instance", handler.Create)
	router.Get("/custom-workout-instance/{id}", handler.GetByID)
	router.Get("/custom-workout-instance", handler.List)
	router.Put("/custom-workout-instance/{id}", handler.Update)
	router.Delete("/custom-workout-instance/{id}", handler.Delete)
	return router
}
