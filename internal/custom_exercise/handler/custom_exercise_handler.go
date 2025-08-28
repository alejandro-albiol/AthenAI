package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomExerciseHandler struct {
	Service interfaces.CustomExerciseService
}

func NewCustomExerciseHandler(svc interfaces.CustomExerciseService) *CustomExerciseHandler {
	return &CustomExerciseHandler{Service: svc}
}

func (h *CustomExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.CustomExerciseCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	gymID := middleware.GetGymID(r)
	err := h.Service.CreateCustomExercise(gymID, dtoReq)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to create custom exercise", err))
		return
	}
	response.WriteAPISuccess(w, "Custom Exercise created successfully", nil)
}

func (h *CustomExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	res, err := h.Service.GetCustomExerciseByID(gymID, id)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve custom exercise", err))
		return
	}
	response.WriteAPISuccess(w, "Custom Exercise retrieved successfully", res)
}

func (h *CustomExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	res, err := h.Service.ListCustomExercises(gymID)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom exercises", err))
		return
	}
	response.WriteAPISuccess(w, "Custom Exercises retrieved successfully", res)
}

func (h *CustomExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.CustomExerciseUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	err := h.Service.UpdateCustomExercise(gymID, id, dtoReq)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to update custom exercise", err))
		return
	}
	response.WriteAPISuccess(w, "Custom Exercise updated successfully", nil)
}

func (h *CustomExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	err := h.Service.DeleteCustomExercise(gymID, id)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom exercise", err))
		return
	}
	response.WriteAPISuccess(w, "Custom Exercise deleted successfully", nil)
}
