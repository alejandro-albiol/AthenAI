package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type GymHandler struct {
	service interfaces.GymService
}

func NewGymHandler(service interfaces.GymService) *GymHandler {
	return &GymHandler{service: service}
}

func (h *GymHandler) CreateGym(w http.ResponseWriter, r *http.Request) {
	// Security validation: only platform admins can create gyms
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can create gyms",
			nil,
		))
		return
	}

	var creationDTO dto.GymCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&creationDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	id, err := h.service.CreateGym(creationDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error when creating gym",
			err,
		))
		return
	}
	response.WriteAPICreated(w, "Gym created successfully", id)
}

func (h *GymHandler) GetGymByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Security validation: ensure user has access to this gym
	if !middleware.ValidateGymAccess(r, id) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: You can only access your own gym data",
			nil,
		))
		return
	}

	gym, err := h.service.GetGymByID(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error when getting gym by id",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Gym found", gym)
}

func (h *GymHandler) GetGymByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	gym, err := h.service.GetGymByName(name)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error when getting gym by domain",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Gym found", gym)
}

func (h *GymHandler) GetAllGyms(w http.ResponseWriter, r *http.Request) {
	// Security validation: only platform admins can get all gyms
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can access all gyms",
			nil,
		))
		return
	}

	gyms, err := h.service.GetAllGyms()
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error when getting all gyms",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Gyms retrieved successfully", gyms)
}

func (h *GymHandler) UpdateGym(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Security validation: ensure user has access to this gym
	if !middleware.ValidateGymAccess(r, id) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: You can only modify your own gym data",
			nil,
		))
		return
	}

	var updateDTO dto.GymUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	updatedGym, err := h.service.UpdateGym(id, updateDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error when updating gym",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Gym updated successfully", updatedGym)
}

func (h *GymHandler) SetGymActive(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	activeStr := r.URL.Query().Get("active")
	active := activeStr == "true"

	err := h.service.SetGymActive(id, active)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error when setting gym active",
			err,
		))
		return
	}
	statusMsg := "deactivated"
	if active {
		statusMsg = "activated"
	}
	response.WriteAPISuccess(w, fmt.Sprintf("Gym %s successfully", statusMsg), nil)
}

func (h *GymHandler) DeleteGym(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.DeleteGym(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error when deleting gym",
			err,
		))
		return
	}
	response.WriteAPISuccess(w, "Gym deleted successfully", nil)
}
