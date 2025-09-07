package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type EquipmentHandler struct {
	service interfaces.EquipmentService
}

func NewEquipmentHandler(service interfaces.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{
		service: service,
	}
}

func (h *EquipmentHandler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	var createDTO *dto.EquipmentCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request body",
			nil,
		))
		return
	}

	equipmentID, err := h.service.CreateEquipment(createDTO)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPICreated(w, "Equipment created successfully", *equipmentID)
}

func (h *EquipmentHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Equipment ID is required",
			nil,
		))
		return
	}

	equipment, err := h.service.GetEquipmentByID(id)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Equipment retrieved successfully", equipment)
}

func (h *EquipmentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	equipment, err := h.service.GetAllEquipment()
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Equipment list retrieved successfully", equipment)
}

func (h *EquipmentHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Equipment ID is required",
			nil,
		))
		return
	}

	var updateDTO *dto.EquipmentUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request body",
			nil,
		))
		return
	}

	equipment, err := h.service.UpdateEquipment(id, updateDTO)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Equipment updated successfully", equipment)
}

func (h *EquipmentHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Equipment ID is required",
			nil,
		))
		return
	}

	err := h.service.DeleteEquipment(id)
	if err != nil {
		response.WriteAPIError(w, err.(*apierror.APIError))
		return
	}

	response.WriteAPISuccess(w, "Equipment deleted successfully", nil)
}
