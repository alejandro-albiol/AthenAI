package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type CustomMemberWorkoutHandler struct {
	Service interfaces.CustomMemberWorkoutService
}

func NewCustomMemberWorkoutHandler(service interfaces.CustomMemberWorkoutService) *CustomMemberWorkoutHandler {
	return &CustomMemberWorkoutHandler{Service: service}
}

func (h *CustomMemberWorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCustomMemberWorkoutDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	gymID := middleware.GetGymID(r)
	id, err := h.Service.CreateCustomMemberWorkout(gymID, &req)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to create custom member workout", err))
		return
	}
	response.WriteAPICreated(w, "Custom member workout created", id)
}

func (h *CustomMemberWorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	res, err := h.Service.GetCustomMemberWorkoutByID(gymID, id)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to get custom member workout", err))
		return
	}
	response.WriteAPISuccess(w, "Custom member workout fetched", res)
}

func (h *CustomMemberWorkoutHandler) ListByMemberID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	memberID := chi.URLParam(r, "memberID")
	res, err := h.Service.ListCustomMemberWorkoutsByMemberID(gymID, memberID)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom member workouts", err))
		return
	}
	response.WriteAPISuccess(w, "Custom member workouts fetched", res)
}

func (h *CustomMemberWorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCustomMemberWorkoutDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	gymID := middleware.GetGymID(r)
	req.ID = chi.URLParam(r, "id")
	err := h.Service.UpdateCustomMemberWorkout(gymID, &req)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to update custom member workout", err))
		return
	}
	response.WriteAPISuccess(w, "Custom member workout updated", nil)
}

func (h *CustomMemberWorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	err := h.Service.DeleteCustomMemberWorkout(gymID, id)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom member workout", err))
		return
	}
	response.WriteAPISuccess(w, "Custom member workout deleted", nil)
}
