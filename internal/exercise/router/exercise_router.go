package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewExerciseRouter(handler interfaces.ExerciseHandler) http.Handler {
	r := chi.NewRouter()

	// Exercise CRUD endpoints
	r.Post("/", handler.CreateExercise)             // POST /exercises
	r.Get("/", handler.GetAllExercises)             // GET /exercises
	r.Get("/search", handler.GetExercisesByFilters) // GET /exercises/search
	r.Get("/{id}", handler.GetExerciseByID)         // GET /exercises/{id}
	r.Put("/{id}", handler.UpdateExercise)          // PUT /exercises/{id}
	r.Delete("/{id}", handler.DeleteExercise)       // DELETE /exercises/{id}

	return r
}
