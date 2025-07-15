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
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type GymHandler struct {
	service interfaces.GymService
}

func NewGymHandler(service interfaces.GymService) *GymHandler {
	return &GymHandler{service: service}
}

func (h *GymHandler) CreateGym(w http.ResponseWriter, r *http.Request) {
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

func (h *GymHandler) GetGymByID(w http.ResponseWriter, r *http.Request, id string) {
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

func (h *GymHandler) GetGymByDomain(w http.ResponseWriter, r *http.Request, domain string) {
	gym, err := h.service.GetGymByDomain(domain)
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

func (h *GymHandler) UpdateGym(w http.ResponseWriter, r *http.Request, id string) {
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

func (h *GymHandler) SetGymActive(w http.ResponseWriter, r *http.Request, id string, active bool) {
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

func (h *GymHandler) DeleteGym(w http.ResponseWriter, r *http.Request, id string) {
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
