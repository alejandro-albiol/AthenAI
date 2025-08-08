package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewExerciseMuscularGroupRouter(handler interfaces.ExerciseMuscularGroupHandler) http.Handler {
	r := chi.NewRouter()

	// Exercise-Muscular Group link endpoints
	r.Post("/link", handler.CreateLink)                                                 // POST /exercise-muscular-groups/link
	r.Get("/link/{id}", handler.GetLinkByID)                                            // GET /exercise-muscular-groups/link/{id}
	r.Delete("/link/{id}", handler.DeleteLink)                                          // DELETE /exercise-muscular-groups/link/{id}
	r.Get("/exercise/{exerciseID}/links", handler.GetLinksByExerciseID)                 // GET /exercise-muscular-groups/exercise/{exerciseID}/links
	r.Get("/muscular-group/{muscularGroupID}/links", handler.GetLinksByMuscularGroupID) // GET /exercise-muscular-groups/muscular-group/{muscularGroupID}/links

	return r
}
