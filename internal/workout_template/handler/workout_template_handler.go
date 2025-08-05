package handler

// WorkoutTemplateHandler handles workout template related requests.
import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/workout_template/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type WorkoutTemplateHandler struct {
	Service interfaces.WorkoutTemplateService
}

// NewWorkoutTemplateHandler creates a new WorkoutTemplateHandler with the provided service.
func NewWorkoutTemplateHandler(service interfaces.WorkoutTemplateService) *WorkoutTemplateHandler {
	return &WorkoutTemplateHandler{
		Service: service,
	}
}

// CreateWorkoutTemplate handles HTTP requests to create a new workout template.
func (h *WorkoutTemplateHandler) CreateWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateWorkoutTemplateDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	err := h.Service.CreateWorkoutTemplate(input)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	response.WriteAPICreated(w, "Workout template created successfully", input)
}

// GetWorkoutTemplateByID handles HTTP requests to retrieve a workout template by ID.
func (h *WorkoutTemplateHandler) GetWorkoutTemplateByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	template, err := h.Service.GetWorkoutTemplateByID(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Workout template found", template)
}

// GetWorkoutTemplateByName handles HTTP requests to retrieve a workout template by name.
func (h *WorkoutTemplateHandler) GetWorkoutTemplateByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	template, err := h.Service.GetWorkoutTemplateByName(name)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Workout template found", template)
}

// GetWorkoutTemplatesByDifficulty handles HTTP requests to retrieve workout templates by difficulty.
func (h *WorkoutTemplateHandler) GetWorkoutTemplatesByDifficulty(w http.ResponseWriter, r *http.Request) {
	difficulty := r.URL.Query().Get("difficulty")
	templates, err := h.Service.GetWorkoutTemplatesByDifficulty(difficulty)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Workout templates retrieved successfully", templates)
}

// GetWorkoutTemplatesByTargetAudience handles HTTP requests to retrieve workout templates by target audience.
func (h *WorkoutTemplateHandler) GetWorkoutTemplatesByTargetAudience(w http.ResponseWriter, r *http.Request) {
	targetAudience := r.URL.Query().Get("target_audience")
	templates, err := h.Service.GetWorkoutTemplatesByTargetAudience(targetAudience)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Workout templates retrieved successfully", templates)
}

// GetAllWorkoutTemplates handles HTTP requests to retrieve all workout templates.
func (h *WorkoutTemplateHandler) GetAllWorkoutTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.Service.GetAllWorkoutTemplates()
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Workout templates retrieved successfully", templates)
}

// UpdateWorkoutTemplate handles HTTP requests to update a workout template.
func (h *WorkoutTemplateHandler) UpdateWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var input dto.UpdateWorkoutTemplateDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	template, err := h.Service.UpdateWorkoutTemplate(id, input)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Workout template updated successfully", template)
}

// DeleteWorkoutTemplate handles HTTP requests to delete a workout template.
func (h *WorkoutTemplateHandler) DeleteWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := h.Service.DeleteWorkoutTemplate(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
