package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}
	id, err := h.service.CreateGym(creationDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			apierror.WriteAPIError(w, apiErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error when creating gym",
			Data:    nil,
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.APIResponse[any]{
		Status:  "success",
		Message: "Gym created successfully",
		Data:    id,
	})
}

func (h *GymHandler) GetGymByID(w http.ResponseWriter, r *http.Request, id string) {
	gym, err := h.service.GetGymByID(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			apierror.WriteAPIError(w, apiErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error when getting gym by id",
			Data:    nil,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[dto.GymResponseDTO]{
		Status:  "success",
		Message: "Gym found",
		Data:    gym,
	})
}

func (h *GymHandler) GetGymByDomain(w http.ResponseWriter, r *http.Request, domain string) {
	gym, err := h.service.GetGymByDomain(domain)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			apierror.WriteAPIError(w, apiErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error when getting gym by domain",
			Data:    nil,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[dto.GymResponseDTO]{
		Status:  "success",
		Message: "Gym found",
		Data:    gym,
	})
}

func (h *GymHandler) GetAllGyms(w http.ResponseWriter, r *http.Request) {
	gyms, err := h.service.GetAllGyms()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error when getting all gyms",
			Data:    nil,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[[]dto.GymResponseDTO]{
		Status:  "success",
		Message: "Gyms retrieved successfully",
		Data:    gyms,
	})
}

func (h *GymHandler) UpdateGym(w http.ResponseWriter, r *http.Request, id string) {
	var updateDTO dto.GymUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}
	updatedGym, err := h.service.UpdateGym(id, updateDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			apierror.WriteAPIError(w, apiErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error when updating gym",
			Data:    nil,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.APIResponse[dto.GymResponseDTO]{
		Status:  "success",
		Message: "Gym updated successfully",
		Data:    updatedGym,
	})
}

func (h *GymHandler) SetGymActive(w http.ResponseWriter, r *http.Request, id string, active bool) {
	err := h.service.SetGymActive(id, active)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			apierror.WriteAPIError(w, apiErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error when setting gym active",
			Data:    nil,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	statusMsg := "deactivated"
	if active {
		statusMsg = "activated"
	}
	json.NewEncoder(w).Encode(response.APIResponse[any]{
		Status:  "success",
		Message: fmt.Sprintf("Gym %s successfully", statusMsg),
		Data:    nil,
	})
}

func (h *GymHandler) DeleteGym(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteGym(id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			apierror.WriteAPIError(w, apiErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.APIResponse[any]{
			Status:  "error",
			Message: "Internal server error when deleting gym",
			Data:    nil,
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
