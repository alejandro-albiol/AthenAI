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
	"github.com/go-chi/chi/v5"
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
	var createDTO dto.CreateWorkoutTemplateDTO
	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	templateID, err := h.Service.CreateWorkoutTemplate(&createDTO)
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
	response.WriteAPICreated(w, "Workout template created successfully", map[string]string{"id": templateID})
}

// GetWorkoutTemplateByID handles HTTP requests to retrieve a workout template by ID.
func (h *WorkoutTemplateHandler) GetWorkoutTemplateByID(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")
	workoutTemplate, err := h.Service.GetWorkoutTemplateByID(templateID)
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
	response.WriteAPISuccess(w, "Workout template found", workoutTemplate)
}

// GetWorkoutTemplateByName handles HTTP requests to retrieve a workout template by name.
func (h *WorkoutTemplateHandler) GetWorkoutTemplateByName(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "name")
	workoutTemplate, err := h.Service.GetWorkoutTemplateByName(templateName)
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
	response.WriteAPISuccess(w, "Workout template found", workoutTemplate)
}

// GetWorkoutTemplatesByDifficulty handles HTTP requests to retrieve workout templates by difficulty.
func (h *WorkoutTemplateHandler) GetWorkoutTemplatesByDifficulty(w http.ResponseWriter, r *http.Request) {
	difficultyLevel := chi.URLParam(r, "difficulty")
	workoutTemplates, err := h.Service.GetWorkoutTemplatesByDifficulty(difficultyLevel)
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
	response.WriteAPISuccess(w, "Workout templates retrieved successfully", workoutTemplates)
}

// GetWorkoutTemplatesByTargetAudience handles HTTP requests to retrieve workout templates by target audience.
func (h *WorkoutTemplateHandler) GetWorkoutTemplatesByTargetAudience(w http.ResponseWriter, r *http.Request) {
	targetAudience := chi.URLParam(r, "targetAudience")
	workoutTemplates, err := h.Service.GetWorkoutTemplatesByTargetAudience(targetAudience)
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
	response.WriteAPISuccess(w, "Workout templates retrieved successfully", workoutTemplates)
}

// GetAllWorkoutTemplates handles HTTP requests to retrieve all workout templates.
func (h *WorkoutTemplateHandler) GetAllWorkoutTemplates(w http.ResponseWriter, r *http.Request) {
	workoutTemplates, err := h.Service.GetAllWorkoutTemplates()
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
	response.WriteAPISuccess(w, "Workout templates retrieved successfully", workoutTemplates)
}

// UpdateWorkoutTemplate handles HTTP requests to update a workout template.
func (h *WorkoutTemplateHandler) UpdateWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")
	var updateDTO dto.UpdateWorkoutTemplateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	workoutTemplate, err := h.Service.UpdateWorkoutTemplate(templateID, &updateDTO)
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
	response.WriteAPISuccess(w, "Workout template updated successfully", workoutTemplate)
}

// DeleteWorkoutTemplate handles HTTP requests to delete a workout template.
func (h *WorkoutTemplateHandler) DeleteWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")
	err := h.Service.DeleteWorkoutTemplate(templateID)
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
