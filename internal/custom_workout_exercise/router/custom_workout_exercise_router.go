package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewCustomWorkoutExerciseRouter(h interfaces.CustomWorkoutExerciseHandler) http.Handler {
	r := chi.NewRouter()

	// POST /custom-workout-exercises - Create a new workout exercise
	r.Post("/custom-workout-exercises", h.Create)

	// GET /custom-workout-exercises/{id} - Get workout exercise by ID
	r.Get("/custom-workout-exercises/{id}", h.GetByID)

	// GET /custom-workout-exercises/workout-instance/{workoutInstanceId} - List exercises by workout instance
	r.Get("/custom-workout-exercises/workout-instance/{workoutInstanceId}", h.ListByWorkoutInstanceID)

	// GET /custom-workout-exercises/muscular-group/{muscularGroupId} - List exercises by muscular group
	r.Get("/custom-workout-exercises/muscular-group/{muscularGroupId}", h.ListByMuscularGroupID)

	// GET /custom-workout-exercises/equipment/{equipmentId} - List exercises by equipment
	r.Get("/custom-workout-exercises/equipment/{equipmentId}", h.ListByEquipmentID)

	// PUT /custom-workout-exercises/{id} - Update workout exercise
	r.Put("/custom-workout-exercises/{id}", h.Update)

	// DELETE /custom-workout-exercises/{id} - Delete workout exercise
	r.Delete("/custom-workout-exercises/{id}", h.Delete)

	return r
}
