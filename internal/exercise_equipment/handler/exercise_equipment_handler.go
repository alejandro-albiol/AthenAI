package handler

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/service"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type ExerciseEquipmentHandler struct {
	service *service.ExerciseEquipmentService
}

func NewExerciseEquipmentHandler(service *service.ExerciseEquipmentService) *ExerciseEquipmentHandler {
	return &ExerciseEquipmentHandler{service: service}
}

// POST /link
func (h *ExerciseEquipmentHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request body to dto.ExerciseEquipment
	// Call service.CreateLink
	response.WriteAPISuccess(w, "Link created", nil)
}

// DELETE /link/{id}
func (h *ExerciseEquipmentHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse id from URL
	// Call service.DeleteLink
	response.WriteAPISuccess(w, "Link deleted", nil)
}

// GET /link/{id}
func (h *ExerciseEquipmentHandler) GetLinkByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse id from URL
	// Call service.GetLinkByID
	response.WriteAPISuccess(w, "Link found", nil)
}

// GET /exercise/{exerciseID}/links
func (h *ExerciseEquipmentHandler) GetLinksByExerciseID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse exerciseID from URL
	// Call service.GetLinksByExerciseID
	response.WriteAPISuccess(w, "Links found", nil)
}

// GET /equipment/{equipmentID}/links
func (h *ExerciseEquipmentHandler) GetLinksByEquipmentID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse equipmentID from URL
	// Call service.GetLinksByEquipmentID
	response.WriteAPISuccess(w, "Links found", nil)
}
