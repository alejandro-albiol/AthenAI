package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_template/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewWorkoutTemplateRouter(handler interfaces.WorkoutTemplateHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateWorkoutTemplate(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAllWorkoutTemplates(w, r)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetWorkoutTemplateByID(w, r)
	})

	r.Get("/name/{name}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetWorkoutTemplateByName(w, r)
	})

	r.Get("difficulty/{difficulty}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetWorkoutTemplatesByDifficulty(w, r)
	})

	r.Get("target-audience/{targetAudience}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetWorkoutTemplatesByTargetAudience(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAllWorkoutTemplates(w, r)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateWorkoutTemplate(w, r)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteWorkoutTemplate(w, r)
	})

	return r
}