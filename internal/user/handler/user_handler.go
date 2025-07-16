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
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type UsersHandler struct {
	service interfaces.UserService
}

func NewUsersHandler(service interfaces.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Ensure user has admin privileges
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can register new users",
			nil,
		))
		return
	}

	var dto dto.UserCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	err := h.service.RegisterUser(gymID, dto)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	response.WriteAPICreated(w, "User registered successfully", dto)
}

func (h *UsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Only admins can list all users
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can list all users",
			nil,
		))
		return
	}

	users, err := h.service.GetAllUsers(gymID)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	response.WriteAPISuccess(w, "Users retrieved successfully", users)
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Users can only access their own profile, admins can access any user in their gym
	currentUserID := middleware.GetUserID(r)
	if !middleware.IsGymAdmin(r) && currentUserID != id {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: You can only access your own profile",
			nil,
		))
		return
	}

	user, err := h.service.GetUserByID(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	response.WriteAPISuccess(w, "User found", user)
}

func (h *UsersHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request, username string) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Only admins can search users by username
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can search users by username",
			nil,
		))
		return
	}

	user, err := h.service.GetUserByUsername(gymID, username)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	response.WriteAPISuccess(w, "User found", user)
}

func (h *UsersHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request, email string) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Only admins can search users by email
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can search users by email",
			nil,
		))
		return
	}

	user, err := h.service.GetUserByEmail(gymID, email)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	response.WriteAPISuccess(w, "User found", user)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id string, userDTO dto.UserUpdateDTO) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Users can only update their own profile, admins can update any user in their gym
	currentUserID := middleware.GetUserID(r)
	if !middleware.IsGymAdmin(r) && currentUserID != id {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: You can only update your own profile",
			nil,
		))
		return
	}

	err := h.service.UpdateUser(gymID, id, userDTO)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	response.WriteAPISuccess(w, "User updated successfully", nil)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Only admins can delete users
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can delete users",
			nil,
		))
		return
	}

	err := h.service.DeleteUser(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UsersHandler) VerifyUser(w http.ResponseWriter, r *http.Request, id string) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Only admins can verify users
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can verify users",
			nil,
		))
		return
	}

	err := h.service.VerifyUser(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	response.WriteAPISuccess(w, "User verified successfully", nil)
}

func (h *UsersHandler) SetUserActive(w http.ResponseWriter, r *http.Request, id string, active bool) {
	// Extract gym ID from JWT token
	gymID := middleware.GetGymID(r)

	// Security: Only admins can activate/deactivate users
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can activate/deactivate users",
			nil,
		))
		return
	}

	err := h.service.SetUserActive(gymID, id, active)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error",
			err,
		))
		return
	}

	statusMsg := "deactivated"
	if active {
		statusMsg = "activated"
	}
	response.WriteAPISuccess(w, fmt.Sprintf("User %s successfully", statusMsg), nil)
}
