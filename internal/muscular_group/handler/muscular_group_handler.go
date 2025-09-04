package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type MuscularGroupHandler struct {
	service interfaces.MuscularGroupService
}

func NewMuscularGroupHandler(service interfaces.MuscularGroupService) *MuscularGroupHandler {
	return &MuscularGroupHandler{
		service: service,
	}
}

func (h *MuscularGroupHandler) CreateMuscularGroup(w http.ResponseWriter, r *http.Request) {
	var createDTO dto.CreateMuscularGroupDTO
	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request body",
			nil,
		))
		return
	}

	muscularGroupID, err := h.service.CreateMuscularGroup(&createDTO)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPICreated(w, "Muscular group created successfully", map[string]string{"id": muscularGroupID})
}

func (h *MuscularGroupHandler) GetMuscularGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Muscular group ID is required",
			nil,
		))
		return
	}

	muscularGroup, err := h.service.GetMuscularGroupByID(id)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Muscular group retrieved successfully", muscularGroup)
}

func (h *MuscularGroupHandler) ListMuscularGroups(w http.ResponseWriter, r *http.Request) {
	muscularGroups, err := h.service.GetAllMuscularGroups()
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Muscular groups retrieved successfully", muscularGroups)
}

func (h *MuscularGroupHandler) UpdateMuscularGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Muscular group ID is required",
			nil,
		))
		return
	}

	var updateDTO *dto.UpdateMuscularGroupDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request body",
			nil,
		))
		return
	}

	muscularGroup, err := h.service.UpdateMuscularGroup(id, updateDTO)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Muscular group updated successfully", muscularGroup)
}

func (h *MuscularGroupHandler) DeleteMuscularGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Muscular group ID is required",
			nil,
		))
		return
	}

	err := h.service.DeleteMuscularGroup(id)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Muscular group deleted successfully", nil)
}
