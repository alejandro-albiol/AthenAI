package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomExerciseMuscularGroupHandler struct {
	service interfaces.CustomExerciseMuscularGroupService
}

func NewCustomExerciseMuscularGroupHandler(svc interfaces.CustomExerciseMuscularGroupService) *CustomExerciseMuscularGroupHandler {
	return &CustomExerciseMuscularGroupHandler{service: svc}
}

func (h *CustomExerciseMuscularGroupHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	var req dto.CustomExerciseMuscularGroupCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	if req.CustomExerciseID == "" || req.MuscularGroupID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing required fields", nil))
		return
	}
	id, err := h.service.CreateLink(gymID, &req)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPICreated(w, "Custom exercise muscular group link created", id)
}

func (h *CustomExerciseMuscularGroupHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	if apiErr, ok := h.service.DeleteLink(gymID, id).(*apierror.APIError); apiErr != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group link deleted", nil)
}

func (h *CustomExerciseMuscularGroupHandler) RemoveAllLinksForExercise(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	customExerciseID := chi.URLParam(r, "customExerciseID")
	if apiErr, ok := h.service.RemoveAllLinksForExercise(gymID, customExerciseID).(*apierror.APIError); apiErr != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "All muscular group links for custom exercise removed", nil)
}

func (h *CustomExerciseMuscularGroupHandler) GetLinkByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	link, err := h.service.GetLinkByID(gymID, id)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group link found", link)
}

func (h *CustomExerciseMuscularGroupHandler) GetLinksByExerciseID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	customExerciseID := chi.URLParam(r, "customExerciseID")
	links, err := h.service.GetLinksByCustomExerciseID(gymID, customExerciseID)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group links for exercise", links)
}

func (h *CustomExerciseMuscularGroupHandler) GetLinksByMuscularGroupID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	muscularGroupID := chi.URLParam(r, "muscularGroupID")
	links, err := h.service.GetLinksByMuscularGroupID(gymID, muscularGroupID)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group links for muscular group", links)
}
