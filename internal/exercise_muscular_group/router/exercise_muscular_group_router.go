package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewExerciseMuscularGroupRouter(handler interfaces.ExerciseMuscularGroupHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/link", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateLink(w, r)
	})
	r.Delete("/link/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteLink(w, r)
	})
	r.Get("/link/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetLinkByID(w, r)
	})
	r.Get("/exercise/{exerciseID}/links", func(w http.ResponseWriter, r *http.Request) {
		handler.GetLinksByExerciseID(w, r)
	})
	r.Get("/muscular-group/{muscularGroupID}/links", func(w http.ResponseWriter, r *http.Request) {
		handler.GetLinksByMuscularGroupID(w, r)
	})

	return r
}
