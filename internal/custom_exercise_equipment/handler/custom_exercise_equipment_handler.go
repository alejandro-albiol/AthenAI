package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomExerciseEquipmentHandler struct {
	service interfaces.CustomExerciseEquipmentService
}

func NewCustomExerciseEquipmentHandler(svc interfaces.CustomExerciseEquipmentService) *CustomExerciseEquipmentHandler {
	return &CustomExerciseEquipmentHandler{service: svc}
}

func (h *CustomExerciseEquipmentHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	var link *dto.CustomExerciseEquipment
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}

	if link.CustomExerciseID == "" || link.EquipmentID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "CustomExerciseID and EquipmentID are required", nil))
		return
	}
	
	id, err := h.service.CreateLink(gymID, link)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise equipment link created", id)
}

func (h *CustomExerciseEquipmentHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	if apiErr, ok := h.service.DeleteLink(gymID, id).(*apierror.APIError); apiErr != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise equipment link deleted", nil)
}

func (h *CustomExerciseEquipmentHandler) RemoveAllLinksForExercise(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	customExerciseID := chi.URLParam(r, "customExerciseID")
	if apiErr, ok := h.service.RemoveAllLinksForExercise(gymID, customExerciseID).(*apierror.APIError); apiErr != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "All equipment links for custom exercise removed", nil)
}

func (h *CustomExerciseEquipmentHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	link, err := h.service.FindByID(gymID, id)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise equipment link found", link)
}

func (h *CustomExerciseEquipmentHandler) FindByCustomExerciseID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	customExerciseID := chi.URLParam(r, "customExerciseID")
	links, err := h.service.FindByCustomExerciseID(gymID, customExerciseID)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise equipment links for exercise", links)
}

func (h *CustomExerciseEquipmentHandler) FindByEquipmentID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	equipmentID := chi.URLParam(r, "equipmentID")
	links, err := h.service.FindByEquipmentID(gymID, equipmentID)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise equipment links for equipment", links)
}
