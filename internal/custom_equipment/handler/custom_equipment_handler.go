package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

// CustomEquipmentHandler handles HTTP requests for custom equipment

type CustomEquipmentHandler struct {
	Service interfaces.CustomEquipmentService
}

func NewCustomEquipmentHandler(service interfaces.CustomEquipmentService) interfaces.CustomEquipmentHandler {
	return &CustomEquipmentHandler{Service: service}
}

func (h *CustomEquipmentHandler) CreateCustomEquipment(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.CreateCustomEquipmentDTO
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	if dtoReq.Name == "" || dtoReq.CreatedBy == "" || dtoReq.Category == "" || dtoReq.Description == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing required fields", nil))
		return
	}
	gymID := middleware.GetGymID(r)
	id, err := h.Service.CreateCustomEquipment(gymID, &dtoReq)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}
	response.WriteAPICreated(w, "Equipment created", id)
}

func (h *CustomEquipmentHandler) GetCustomEquipmentByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := r.URL.Query().Get("id")
	equipment, err := h.Service.GetCustomEquipmentByID(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Equipment found", equipment)
}

// GetCustomEquipmentByName handles GET /custom-equipment/search?name=...
func (h *CustomEquipmentHandler) GetCustomEquipmentByName(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	name := r.URL.Query().Get("name")
	equipment, err := h.Service.GetCustomEquipmentByName(gymID, name)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Equipment found", equipment)
}

func (h *CustomEquipmentHandler) ListCustomEquipment(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	equipment, err := h.Service.ListCustomEquipment(gymID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Equipment list", equipment)
}

func (h *CustomEquipmentHandler) UpdateCustomEquipment(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.UpdateCustomEquipmentDTO
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	gymID := middleware.GetGymID(r)
	err := h.Service.UpdateCustomEquipment(gymID, &dtoReq)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Equipment updated", nil)
}

func (h *CustomEquipmentHandler) DeleteCustomEquipment(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := r.URL.Query().Get("id")
	err := h.Service.DeleteCustomEquipment(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unexpected error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Equipment deleted", nil)
}
