package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomWorkoutExerciseHandler struct {
	Service interfaces.CustomWorkoutExerciseService
}

func NewCustomWorkoutExerciseHandler(service interfaces.CustomWorkoutExerciseService) *CustomWorkoutExerciseHandler {
	return &CustomWorkoutExerciseHandler{
		Service: service,
	}
}

func (h *CustomWorkoutExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.CreateCustomWorkoutExerciseDTO
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}

	gymID := middleware.GetGymID(r)

	id, err := h.Service.CreateCustomWorkoutExercise(gymID, &dtoReq)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}

	response.WriteAPICreated(w, "Workout exercise created successfully", id)
}

func (h *CustomWorkoutExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")

	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Exercise ID is required", nil))
		return
	}

	exercise, err := h.Service.GetCustomWorkoutExerciseByID(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}

	response.WriteAPISuccess(w, "Workout exercise retrieved successfully", exercise)
}

func (h *CustomWorkoutExerciseHandler) ListByWorkoutInstanceID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	workoutInstanceID := chi.URLParam(r, "workoutInstanceId")

	if workoutInstanceID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Workout instance ID is required", nil))
		return
	}

	exercises, err := h.Service.ListCustomWorkoutExercisesByWorkoutInstanceID(gymID, workoutInstanceID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}

	response.WriteAPISuccess(w, "Workout exercises retrieved successfully", exercises)
}

func (h *CustomWorkoutExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.UpdateCustomWorkoutExerciseDTO
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}

	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")

	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Exercise ID is required", nil))
		return
	}

	// Set the ID from URL parameter
	dtoReq.ID = id

	err := h.Service.UpdateCustomWorkoutExercise(gymID, &dtoReq)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}

	response.WriteAPISuccess(w, "Workout exercise updated successfully", nil)
}

func (h *CustomWorkoutExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")

	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Exercise ID is required", nil))
		return
	}

	err := h.Service.DeleteCustomWorkoutExercise(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}

	response.WriteAPISuccess(w, "Workout exercise deleted successfully", nil)
}

func (h *CustomWorkoutExerciseHandler) ListByMuscularGroupID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	muscularGroupID := chi.URLParam(r, "muscularGroupId")

	if muscularGroupID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Muscular group ID is required", nil))
		return
	}

	exercises, err := h.Service.ListCustomWorkoutExercisesByMuscularGroupID(gymID, muscularGroupID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}

	response.WriteAPISuccess(w, "Workout exercises retrieved successfully", exercises)
}

func (h *CustomWorkoutExerciseHandler) ListByEquipmentID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	equipmentID := chi.URLParam(r, "equipmentId")

	if equipmentID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Equipment ID is required", nil))
		return
	}

	exercises, err := h.Service.ListCustomWorkoutExercisesByEquipmentID(gymID, equipmentID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}

	response.WriteAPISuccess(w, "Workout exercises retrieved successfully", exercises)
}
