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
	var link dto.CustomExerciseMuscularGroup
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	err := h.service.CreateLink(gymID, link)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group link created", nil)
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

func (h *CustomExerciseMuscularGroupHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	link, err := h.service.GetLinkByID(gymID, id)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group link found", link)
}

func (h *CustomExerciseMuscularGroupHandler) FindByCustomExerciseID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	customExerciseID := chi.URLParam(r, "customExerciseID")
	links, err := h.service.GetLinksByCustomExerciseID(gymID, customExerciseID)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group links for exercise", links)
}

func (h *CustomExerciseMuscularGroupHandler) FindByMuscularGroupID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	muscularGroupID := chi.URLParam(r, "muscularGroupID")
	links, err := h.service.GetLinksByMuscularGroupID(gymID, muscularGroupID)
	if apiErr, ok := err.(*apierror.APIError); err != nil && ok {
		response.WriteAPIError(w, apiErr)
		return
	}
	response.WriteAPISuccess(w, "Custom exercise muscular group links for muscular group", links)
}
