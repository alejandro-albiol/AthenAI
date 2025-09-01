package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomWorkoutInstanceHandler struct {
	Service interfaces.CustomWorkoutInstanceService
}

func NewCustomWorkoutInstanceHandler(service interfaces.CustomWorkoutInstanceService) *CustomWorkoutInstanceHandler {
	return &CustomWorkoutInstanceHandler{Service: service}
}

// API: POST /custom-workout-instance
func (h *CustomWorkoutInstanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	var reqBody dto.CreateCustomWorkoutInstanceDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Malformed JSON", err))
		return
	}
	id, err := h.Service.CreateCustomWorkoutInstance(gymID, reqBody)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPICreated(w, "Workout instance created", id)
}

// API: GET /custom-workout-instance/{id}
func (h *CustomWorkoutInstanceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	id := chi.URLParam(r, "id")
	instance, err := h.Service.GetCustomWorkoutInstanceByID(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Workout instance found", instance)
}

// API: GET /custom-workout-instance
func (h *CustomWorkoutInstanceHandler) List(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	instances, err := h.Service.ListCustomWorkoutInstances(gymID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Workout instances listed", instances)
}

// API: PUT /custom-workout-instance/{id}
func (h *CustomWorkoutInstanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	id := chi.URLParam(r, "id")
	var reqBody dto.UpdateCustomWorkoutInstanceDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Malformed JSON", err))
		return
	}
	reqBody.ID = id
	err = h.Service.UpdateCustomWorkoutInstance(gymID, reqBody)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Workout instance updated", nil)
}

// API: DELETE /custom-workout-instance/{id}
func (h *CustomWorkoutInstanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	id := chi.URLParam(r, "id")
	err := h.Service.DeleteCustomWorkoutInstance(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Workout instance deleted", nil)
}
