package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type TemplateBlockHandler struct {
	service interfaces.TemplateBlockService
}

func NewTemplateBlockHandler(service interfaces.TemplateBlockService) *TemplateBlockHandler {
	return &TemplateBlockHandler{service: service}
}

func (h *TemplateBlockHandler) CreateTemplateBlock(w http.ResponseWriter, r *http.Request) {
	var block dto.CreateTemplateBlockDTO
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	// Validate required fields
	if block.TemplateID == "" || block.BlockName == "" || block.BlockType == "" || block.BlockOrder == 0 || block.ExerciseCount == 0 {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing or invalid required fields: template_id, block_name, block_type, block_order, exercise_count", nil))
		return
	}
	blockID, err := h.service.CreateTemplateBlock(&block)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to create template block", err))
		}
		return
	}
	response.WriteAPICreated(w, "Template block created successfully", *blockID)
}

func (h *TemplateBlockHandler) GetTemplateBlockByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing block ID", nil))
		return
	}
	block, err := h.service.GetTemplateBlockByID(id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to get template block", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Template block retrieved successfully", block)
}

func (h *TemplateBlockHandler) UpdateTemplateBlock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing block ID", nil))
		return
	}
	var update *dto.UpdateTemplateBlockDTO
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	updatedBlock, err := h.service.UpdateTemplateBlock(id, update)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to update template block", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Template block updated successfully", updatedBlock)
}

func (h *TemplateBlockHandler) DeleteTemplateBlock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing block ID", nil))
		return
	}
	if err := h.service.DeleteTemplateBlock(id); err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to delete template block", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Template block deleted successfully", nil)
}

func (h *TemplateBlockHandler) ListTemplateBlocksByTemplateID(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "templateId")
	if templateID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing template ID", nil))
		return
	}
	blocks, err := h.service.ListTemplateBlocksByTemplateID(templateID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to list template blocks", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Template blocks retrieved successfully", blocks)
}
