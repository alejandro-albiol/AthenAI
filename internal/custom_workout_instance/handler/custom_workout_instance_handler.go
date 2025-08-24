package handler

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/interfaces"
)

type CustomWorkoutInstanceHandler struct {
	Service interfaces.CustomWorkoutInstanceService
}

func NewCustomWorkoutInstanceHandler(service interfaces.CustomWorkoutInstanceService) *CustomWorkoutInstanceHandler {
	return &CustomWorkoutInstanceHandler{Service: service}
}

// API: POST /custom-workout-instance
func (h *CustomWorkoutInstanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	// response.WriteAPISuccess(w, nil)
}

// API: GET /custom-workout-instance/{id}
func (h *CustomWorkoutInstanceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	// response.WriteAPISuccess(w, nil)
}

// API: GET /custom-workout-instance
func (h *CustomWorkoutInstanceHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	// response.WriteAPISuccess(w, nil)
}

// API: PUT /custom-workout-instance/{id}
func (h *CustomWorkoutInstanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	// response.WriteAPISuccess(w, nil)
}

// API: DELETE /custom-workout-instance/{id}
func (h *CustomWorkoutInstanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	// response.WriteAPISuccess(w, nil)
}
