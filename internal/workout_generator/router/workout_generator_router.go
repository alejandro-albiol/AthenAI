package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_generator/interfaces"
	"github.com/go-chi/chi/v5"
)

// NewWorkoutGeneratorRouter mounts workout generator endpoints with gym context
func NewWorkoutGeneratorRouter(handler interfaces.WorkoutGeneratorHandler) http.Handler {

	r := chi.NewRouter()
	r.Post("/generate", handler.GenerateWorkout)
	return r
}
