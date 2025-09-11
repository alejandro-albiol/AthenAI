package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewCustomWorkoutInstanceRouter(handler interfaces.CustomWorkoutInstanceHandler) http.Handler {
	router := chi.NewRouter()

	// Basic CRUD operations
	router.Post("/custom-workout-instance", handler.Create)
	router.Get("/custom-workout-instance/{id}", handler.GetByID)
	router.Get("/custom-workout-instance", handler.List)
	router.Put("/custom-workout-instance/{id}", handler.Update)
	router.Delete("/custom-workout-instance/{id}", handler.Delete)

	// Summary operations
	router.Get("/custom-workout-instance/{id}/summary", handler.GetSummaryByID)
	router.Get("/custom-workout-instance/summaries", handler.ListSummaries)

	// User-specific operations
	router.Get("/custom-workout-instance/user/{userID}", handler.GetByUserID)
	router.Get("/custom-workout-instance/user/{userID}/summaries", handler.GetSummariesByUserID)
	router.Get("/custom-workout-instance/user/{userID}/last", handler.GetLastsByUserID)

	return router
}
