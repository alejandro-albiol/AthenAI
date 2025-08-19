package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_generator/dto"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

// WorkoutGeneratorHandler implements interfaces.WorkoutGeneratorHandler
type WorkoutGeneratorHandler struct {
	service interfaces.WorkoutGeneratorService
}

// NewWorkoutGeneratorHandler creates a new WorkoutGeneratorHandler
func NewWorkoutGeneratorHandler(service interfaces.WorkoutGeneratorService) *WorkoutGeneratorHandler {
	return &WorkoutGeneratorHandler{service: service}
}

func (h *WorkoutGeneratorHandler) GenerateWorkout(w http.ResponseWriter, r *http.Request) {
	var req dto.WorkoutGeneratorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "invalid request", err))
		return
	}

	resp, err := h.service.GenerateWorkout(req)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "unexpected error", err))
		}
		return
	}

	response.WriteAPISuccess(w, "workout generated", resp)
}
