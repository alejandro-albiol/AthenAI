package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomExerciseHandler struct {
	Service *service.CustomExerciseService
}

func NewCustomExerciseHandler(svc *service.CustomExerciseService) *CustomExerciseHandler {
	return &CustomExerciseHandler{Service: svc}
}

func (h *CustomExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.CustomExerciseCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	gymID := middleware.GetGymID(r)
	err := h.Service.Create(gymID, &dtoReq)
	if err != nil {
		response.WriteAPIError(w, err)
		return
	}
	response.WriteAPISuccess(w, "Equipment created", nil)
}

func (h *CustomExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	res, err := h.Service.GetByID(gymID, id)
	if err != nil {
		response.WriteAPIError(w, err)
		return
	}
	response.WriteAPISuccess(w, "Custom Exercise retrieved successfully", res)
}

func (h *CustomExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	res, err := h.Service.List(gymID)
	if err != nil {
		response.WriteAPIError(w, err)
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
	dtoReq.ID = &id
	err := h.Service.Update(gymID, &dtoReq)
	if err != nil {
		response.WriteAPIError(w, err)
		return
	}
	response.WriteAPISuccess(w, "Custom Exercise updated successfully", nil)
}

func (h *CustomExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	err := h.Service.Delete(gymID, id)
	if err != nil {
		response.WriteAPIError(w, err)
		return
	}
	response.WriteAPISuccess(w, "Custom Exercise deleted successfully", nil)
}
