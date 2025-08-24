package handler

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type CustomMemberWorkoutHandler struct {
	Service interfaces.CustomMemberWorkoutService
}

func NewCustomMemberWorkoutHandler(service interfaces.CustomMemberWorkoutService) *CustomMemberWorkoutHandler {
	return &CustomMemberWorkoutHandler{Service: service}
}

func (h *CustomMemberWorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomMemberWorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomMemberWorkoutHandler) ListByMemberID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomMemberWorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomMemberWorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}
