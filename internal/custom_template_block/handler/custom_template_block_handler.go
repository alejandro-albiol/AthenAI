package handler

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/service"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type CustomTemplateBlockHandler struct {
	Service *service.CustomTemplateBlockServiceImpl
}

func NewCustomTemplateBlockHandler(service *service.CustomTemplateBlockServiceImpl) *CustomTemplateBlockHandler {
	return &CustomTemplateBlockHandler{Service: service}
}

func (h *CustomTemplateBlockHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomTemplateBlockHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomTemplateBlockHandler) ListByTemplateID(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomTemplateBlockHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}

func (h *CustomTemplateBlockHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call service, write response
	response.WriteAPISuccess(w, "", nil)
}
