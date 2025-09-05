package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	"github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type ExerciseEquipmentHandler struct {
	service interfaces.ExerciseEquipmentService
}

func NewExerciseEquipmentHandler(service interfaces.ExerciseEquipmentService) *ExerciseEquipmentHandler {
	return &ExerciseEquipmentHandler{service: service}
}

// POST /link
func (h *ExerciseEquipmentHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link dto.ExerciseEquipment
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	linkID, err := h.service.CreateLink(&link)
	if err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to create link",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Link created", linkID)
}

// DELETE /link/{id}
func (h *ExerciseEquipmentHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteLink(id); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to delete link",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Link deleted", nil)
}

// GET /link/{id}
func (h *ExerciseEquipmentHandler) GetLinkByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	link, err := h.service.GetLinkByID(id)
	if err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to get link",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Link found", link)
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
