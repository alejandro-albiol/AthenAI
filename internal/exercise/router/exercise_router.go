package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewExerciseRouter(handler interfaces.ExerciseHandler) http.Handler {

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAllExercises(w, r)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateExercise(w, r)
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetExerciseByID(w, r, chi.URLParam(r, "id"))
	})
	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateExercise(w, r, chi.URLParam(r, "id"))
	})
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteExercise(w, r, chi.URLParam(r, "id"))
	})
	return r
}
