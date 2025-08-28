package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/handler"
	"github.com/go-chi/chi/v5"
)

func NewCustomExerciseEquipmentRouter(h *handler.CustomExerciseEquipmentHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/custom-exercise-equipment/link", h.CreateLink)
	r.Delete("/custom-exercise-equipment/link/{id}", h.DeleteLink)
	r.Delete("/custom-exercise-equipment/exercise/{customExerciseID}/links", h.RemoveAllLinksForExercise)
	r.Get("/custom-exercise-equipment/link/{id}", h.FindByID)
	r.Get("/custom-exercise-equipment/exercise/{customExerciseID}/links", h.FindByCustomExerciseID)
	r.Get("/custom-exercise-equipment/equipment/{equipmentID}/links", h.FindByEquipmentID)

	return r
}
