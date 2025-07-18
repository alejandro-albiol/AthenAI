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

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		groups := r.URL.Query()["group"]
		equipment := r.URL.Query()["equipment"]
		handler.GetExercisesByFilters(w, r, groups, equipment)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.UpdateExercise(w, r, id)
	})
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.DeleteExercise(w, r, id)
	})
	return r
}
