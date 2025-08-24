package handler

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/service"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type CustomWorkoutExerciseHandler struct {
	Service *service.CustomWorkoutExerciseService
}

func NewCustomWorkoutExerciseHandler(service *service.CustomWorkoutExerciseService) *CustomWorkoutExerciseHandler {
	return &CustomWorkoutExerciseHandler{
		Service: service,
	}
}

func (h *CustomWorkoutExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutExerciseHandler) ListByWorkoutInstanceID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}
