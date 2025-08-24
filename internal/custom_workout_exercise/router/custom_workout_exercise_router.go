package router

import (
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/handler"
	"github.com/go-chi/chi/v5"
)

func NewCustomWorkoutExerciseRouter(h *handler.CustomWorkoutExerciseHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/custom-workout-exercise", h.Create)
	r.Get("/custom-workout-exercise/{id}", h.GetByID)
	r.Get("/custom-workout-exercise/workout/{workoutInstanceID}", h.ListByWorkoutInstanceID)
	r.Put("/custom-workout-exercise/{id}", h.Update)
	r.Delete("/custom-workout-exercise/{id}", h.Delete)
	return r
}
