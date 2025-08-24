package handler

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/service"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type CustomWorkoutTemplateHandler struct {
	Service *service.CustomWorkoutTemplateServiceImpl
}

func NewCustomWorkoutTemplateHandler(service *service.CustomWorkoutTemplateServiceImpl) *CustomWorkoutTemplateHandler {
	return &CustomWorkoutTemplateHandler{Service: service}
}

func (h *CustomWorkoutTemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutTemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutTemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutTemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomWorkoutTemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}
