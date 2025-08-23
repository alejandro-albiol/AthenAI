package handler

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/service"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

// CustomEquipmentHandler handles HTTP requests for custom equipment

type CustomEquipmentHandler struct {
	Service *service.CustomEquipmentServiceImpl
}

func (h *CustomEquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomEquipmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomEquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomEquipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomEquipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}
