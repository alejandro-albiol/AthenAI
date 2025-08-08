package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type ExerciseMuscularGroupHandler struct {
	service *service.ExerciseMuscularGroupService
}

func NewExerciseMuscularGroupHandler(service *service.ExerciseMuscularGroupService) *ExerciseMuscularGroupHandler {
	return &ExerciseMuscularGroupHandler{service: service}
}

func (h *ExerciseMuscularGroupHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link dto.ExerciseMuscularGroup
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	linkID, err := h.service.CreateLink(link)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to create link",
				err,
			))
		}
		return
	}
	response.WriteAPICreated(w, "Link created successfully", map[string]string{"id": linkID})
}

func (h *ExerciseMuscularGroupHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteLink(id); err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to delete link",
				err,
			))
		}
		return
	}
	response.WriteAPISuccess(w, "Link deleted successfully", nil)
}

func (h *ExerciseMuscularGroupHandler) GetLinkByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	link, err := h.service.GetLinkByID(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to retrieve link",
				err,
			))
		}
		return
	}
	response.WriteAPISuccess(w, "Link retrieved successfully", link)
}

func (h *ExerciseMuscularGroupHandler) GetLinksByExerciseID(w http.ResponseWriter, r *http.Request) {
	exerciseID := chi.URLParam(r, "exerciseID")
	links, err := h.service.GetLinksByExerciseID(exerciseID)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to retrieve links",
				err,
			))
		}
		return
	}
	response.WriteAPISuccess(w, "Links retrieved successfully", links)
}

func (h *ExerciseMuscularGroupHandler) GetLinksByMuscularGroupID(w http.ResponseWriter, r *http.Request) {
	muscularGroupID := chi.URLParam(r, "muscularGroupID")
	links, err := h.service.GetLinksByMuscularGroupID(muscularGroupID)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to retrieve links",
				err,
			))
		}
		return
	}
	response.WriteAPISuccess(w, "Links retrieved successfully", links)
}
