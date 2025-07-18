package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/handler"
	"github.com/go-chi/chi/v5"
)

func NewExerciseEquipmentRouter(h *handler.ExerciseEquipmentHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/link", h.CreateLink)
	r.Delete("/link/{id}", h.DeleteLink)
	r.Get("/link/{id}", h.GetLinkByID)
	r.Get("/exercise/{exerciseID}/links", h.GetLinksByExerciseID)
	r.Get("/equipment/{equipmentID}/links", h.GetLinksByEquipmentID)

	return r
}
