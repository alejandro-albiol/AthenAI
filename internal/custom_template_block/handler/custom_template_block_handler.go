package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomTemplateBlockHandler struct {
	Service interfaces.CustomTemplateBlockService
}

func NewCustomTemplateBlockHandler(service interfaces.CustomTemplateBlockService) *CustomTemplateBlockHandler {
	return &CustomTemplateBlockHandler{Service: service}
}

func (h *CustomTemplateBlockHandler) Create(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing gym ID", nil))
		return
	}

	var block dto.CreateCustomTemplateBlockDTO
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}

	// Validate required fields
	if block.TemplateID == "" || block.BlockName == "" || block.BlockType == "" || block.BlockOrder == 0 || block.ExerciseCount == 0 || block.CreatedBy == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing or invalid required fields: template_id, block_name, block_type, block_order, exercise_count, created_by", nil))
		return
	}

	blockID, err := h.Service.CreateCustomTemplateBlock(gymID, &block)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to create custom template block", err))
		}
		return
	}
	response.WriteAPICreated(w, "Custom template block created successfully", *blockID)
}

func (h *CustomTemplateBlockHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing gym ID", nil))
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing block ID", nil))
		return
	}

	block, err := h.Service.GetCustomTemplateBlockByID(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to get custom template block", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom template block retrieved successfully", block)
}

func (h *CustomTemplateBlockHandler) ListByTemplateID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing gym ID", nil))
		return
	}

	templateID := chi.URLParam(r, "templateId")
	if templateID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing template ID", nil))
		return
	}

	blocks, err := h.Service.ListCustomTemplateBlocksByTemplateID(gymID, templateID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom template blocks", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom template blocks retrieved successfully", blocks)
}

func (h *CustomTemplateBlockHandler) Update(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing gym ID", nil))
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing block ID", nil))
		return
	}

	var update dto.UpdateCustomTemplateBlockDTO
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}

	updatedBlock, err := h.Service.UpdateCustomTemplateBlock(gymID, id, &update)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to update custom template block", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom template block updated successfully", updatedBlock)
}

func (h *CustomTemplateBlockHandler) Delete(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing gym ID", nil))
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Missing block ID", nil))
		return
	}

	if err := h.Service.DeleteCustomTemplateBlock(gymID, id); err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom template block", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom template block deleted successfully", nil)
}
