package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewCustomExerciseMuscularGroupRouter(h interfaces.CustomExerciseMuscularGroupHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/custom_exercise_muscular_group", h.CreateLink)
	r.Get("/custom_exercise_muscular_group/{id}", h.GetLinkByID)
	r.Get("/custom_exercise_muscular_group/custom_exercise/{customExerciseID}", h.GetLinksByExerciseID)
	r.Get("/custom_exercise_muscular_group/muscular_group/{muscularGroupID}", h.GetLinksByMuscularGroupID)
	r.Delete("/custom_exercise_muscular_group/remove_all/custom_exercise/{customExerciseID}", h.RemoveAllLinksForExercise)
	r.Delete("/custom_exercise_muscular_group/{id}", h.DeleteLink)
	return r
}
