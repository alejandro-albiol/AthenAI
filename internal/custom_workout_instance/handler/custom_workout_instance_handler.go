package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
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

	// Get user ID from JWT token
	createdBy := middleware.GetUserID(r)
	if createdBy == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeUnauthorized, "User ID not found in token", nil))
		return
	}

	id, err := h.Service.CreateCustomWorkoutInstance(gymID, createdBy, &reqBody)
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

// API: GET /custom-workout-instance/{id}/summary
func (h *CustomWorkoutInstanceHandler) GetSummaryByID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	id := chi.URLParam(r, "id")
	instance, err := h.Service.GetCustomWorkoutInstanceSummaryByID(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Workout instance summary found", instance)
}

// API: GET /custom-workout-instance/user/{userID}
func (h *CustomWorkoutInstanceHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	userID := chi.URLParam(r, "userID")
	instances, err := h.Service.GetCustomWorkoutInstancesByUserID(gymID, userID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "User workout instances found", instances)
}

// API: GET /custom-workout-instance/user/{userID}/summaries
func (h *CustomWorkoutInstanceHandler) GetSummariesByUserID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	userID := chi.URLParam(r, "userID")
	instances, err := h.Service.GetCustomWorkoutInstanceSummariesByUserID(gymID, userID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "User workout instance summaries found", instances)
}

// API: GET /custom-workout-instance/user/{userID}/last?count={count}
func (h *CustomWorkoutInstanceHandler) GetLastsByUserID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	userID := chi.URLParam(r, "userID")

	// Get count parameter from query string, default to 10
	countStr := r.URL.Query().Get("count")
	count := 10 // default
	if countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil && c > 0 {
			count = c
		}
	}

	instances, err := h.Service.GetLastCustomWorkoutInstancesByUserID(gymID, userID, count)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Last workout instances found", instances)
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

// API: GET /custom-workout-instance/summaries
func (h *CustomWorkoutInstanceHandler) ListSummaries(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	instances, err := h.Service.ListCustomWorkoutInstanceSummaries(gymID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Workout instance summaries listed", instances)
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

	err = h.Service.UpdateCustomWorkoutInstance(gymID, id, &reqBody)
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
