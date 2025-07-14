package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type UsersHandler struct {
	service interfaces.UserService
}

func NewUsersHandler(service interfaces.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request, gymID string) {
	var dto dto.UserCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		apiErr := apierror.New(errorcode_enum.CodeBadRequest, "Invalid request payload", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	err := h.service.RegisterUser(gymID, dto)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "User registered successfully", dto)
}

func (h *UsersHandler) GetAllUsers(w http.ResponseWriter, gymID string) {
	users, err := h.service.GetAllUsers(gymID)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Users retrieved successfully", users)
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, gymID, id string) {
	user, err := h.service.GetUserByID(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "User found", user)
}

func (h *UsersHandler) GetUserByUsername(w http.ResponseWriter, gymID, username string) {
	user, err := h.service.GetUserByUsername(gymID, username)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "User found", user)
}

func (h *UsersHandler) GetUserByEmail(w http.ResponseWriter, gymID, email string) {
	user, err := h.service.GetUserByEmail(gymID, email)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "User found", user)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, gymID, id string, userDTO dto.UserUpdateDTO) {
	err := h.service.UpdateUser(gymID, id, userDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "User updated successfully", nil)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, gymID, id string) {
	err := h.service.DeleteUser(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UsersHandler) VerifyUser(w http.ResponseWriter, gymID, id string) {
	err := h.service.VerifyUser(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "User verified successfully", nil)
}

func (h *UsersHandler) SetUserActive(w http.ResponseWriter, gymID, id string, active bool) {
	err := h.service.SetUserActive(gymID, id, active)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		apiErr = apierror.New(errorcode_enum.CodeInternal, "Internal server error", err)
		response.WriteAPIError(w, apiErr)
		return
	}

	statusMsg := "deactivated"
	if active {
		statusMsg = "activated"
	}
	response.WriteAPISuccess(w, fmt.Sprintf("User %s successfully", statusMsg), nil)
}
