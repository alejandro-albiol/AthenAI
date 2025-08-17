package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_generator/dto"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/interfaces"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	resp, err := h.service.GenerateWorkout(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
