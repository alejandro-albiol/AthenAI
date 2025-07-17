package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type ExerciseHandler struct {
	service *service.ExerciseService
}

func NewExerciseHandler(service *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{service: service}
}

func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {

	creationDTO := dto.ExerciseCreationDTO{}
	if err := json.NewDecoder(r.Body).Decode(&creationDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	id, err := h.service.CreateExercise(creationDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to create exercise",
				err,
			))
		}
		return
	}
	response.WriteAPICreated(w, "Exercise created successfully", id)
}

func (h *ExerciseHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request, id string) {
	exercise, err := h.service.GetExerciseByID(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to retrieve exercise",
				err,
			))
		}
		return
	}
	response.WriteAPISuccess(w, "Exercise retrieved successfully", exercise)
}

func (h *ExerciseHandler) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.service.GetAllExercises()
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to retrieve exercises",
				err,
			))
		}
		return
	}
	response.WriteAPISuccess(w, "Exercises retrieved successfully", exercises)
}

func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request, id string) {
	updateDTO := dto.ExerciseUpdateDTO{}
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	updatedExercise, err := h.service.UpdateExercise(id, updateDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to update exercise",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Exercise updated successfully", updatedExercise)
}

func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteExercise(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to delete exercise",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Exercise deleted successfully", nil)
}