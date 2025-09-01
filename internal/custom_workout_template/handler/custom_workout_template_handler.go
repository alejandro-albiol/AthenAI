package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CustomWorkoutTemplateHandler struct {
	Service *service.CustomWorkoutTemplateServiceImpl
}

func NewCustomWorkoutTemplateHandler(service *service.CustomWorkoutTemplateServiceImpl) *CustomWorkoutTemplateHandler {
	return &CustomWorkoutTemplateHandler{Service: service}
}

func (h *CustomWorkoutTemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	var reqBody dto.CreateCustomWorkoutTemplateDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Malformed JSON", err))
		return
	}
	id, err := h.Service.CreateCustomWorkoutTemplate(gymID, reqBody)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPICreated(w, "Custom workout template created", id)
}

func (h *CustomWorkoutTemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	id := chi.URLParam(r, "id")
	template, err := h.Service.GetCustomWorkoutTemplateByID(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom workout template found", template)
}

func (h *CustomWorkoutTemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	templates, err := h.Service.ListCustomWorkoutTemplates(gymID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom workout templates listed", templates)
}

func (h *CustomWorkoutTemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	id := chi.URLParam(r, "id")
	var reqBody dto.UpdateCustomWorkoutTemplateDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Invalid request body", err))
		return
	}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		response.WriteAPIError(w, apierror.New(errorcode_enum.CodeBadRequest, "Malformed JSON", err))
		return
	}
	reqBody.ID = id
	err = h.Service.UpdateCustomWorkoutTemplate(gymID, reqBody)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom workout template updated", nil)
}

func (h *CustomWorkoutTemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	gymID := chi.URLParam(r, "gymID")
	id := chi.URLParam(r, "id")
	err := h.Service.DeleteCustomWorkoutTemplate(gymID, id)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.WriteAPIError(w, apiErr)
		} else {
			response.WriteAPIError(w, apierror.New(errorcode_enum.CodeInternal, "Unknown error", err))
		}
		return
	}
	response.WriteAPISuccess(w, "Custom workout template deleted", nil)
}
