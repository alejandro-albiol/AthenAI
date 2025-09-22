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
	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	service interfaces.UserService
}

func NewUsersHandler(service interfaces.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Extract gym ID from JWT token or header
	gymID := middleware.GetGymID(r)
	requesterRole := middleware.GetUserRole(r)

	var userDTO dto.UserCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	// Validate required fields (username, email, password, role)
	if userDTO.Username == "" || userDTO.Email == "" || userDTO.Password == "" || userDTO.Role == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Missing required user fields",
			nil,
		))
		return
	}

	// Registration logic by requester role
	switch requesterRole {
	case "": // Public/self-registration
		// Only allow self-registration as user, gymID must be present (from header or JWT)
		if gymID == "" {
			gymID = r.Header.Get("X-Gym-ID")
		}
		if gymID == "" {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeBadRequest,
				"Missing or invalid gym ID for self-registration",
				nil,
			))
			return
		}
		if userDTO.Role != "user" {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeForbidden,
				"Only user role can self-register",
				nil,
			))
			return
		}
	case "superadmin":
		// Can register tenant admins or users for any gym (gymID required)
		if gymID == "" {
			gymID = r.Header.Get("X-Gym-ID")
		}
		if gymID == "" {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeBadRequest,
				"Missing or invalid gym ID for superadmin registration",
				nil,
			))
			return
		}
		if userDTO.Role != "admin" && userDTO.Role != "user" {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeForbidden,
				"Superadmin can only register admin or user roles",
				nil,
			))
			return
		}
	case "admin":
		// Can only register users for their own gym
		if gymID == "" {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeBadRequest,
				"Missing or invalid gym ID for tenant admin registration",
				nil,
			))
			return
		}
		if userDTO.Role != "user" {
			response.WriteAPIError(w, apierror.New(
				errorcode_enum.CodeForbidden,
				"Tenant admin can only register user roles",
				nil,
			))
			return
		}
	default:
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: insufficient privileges to register user",
			nil,
		))
		return
	}

	userID, err := h.service.RegisterUser(gymID, &userDTO)
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

	response.WriteAPICreated(w, "User registered successfully", *userID)
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

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
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

func (h *UsersHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	username := chi.URLParam(r, "username")
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

func (h *UsersHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	email := chi.URLParam(r, "email")
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

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	currentUserID := middleware.GetUserID(r)
	if !middleware.IsGymAdmin(r) && currentUserID != id {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: You can only update your own profile",
			nil,
		))
		return
	}
	userDTO := dto.UserUpdateDTO{}
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	err := h.service.UpdateUser(gymID, id, &userDTO)
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

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
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

func (h *UsersHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
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

func (h *UsersHandler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	gymID := middleware.GetGymID(r)
	id := chi.URLParam(r, "id")
	var activeReq struct{ Active bool }
	if err := json.NewDecoder(r.Body).Decode(&activeReq); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}
	if !middleware.IsGymAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only administrators can activate/deactivate users",
			nil,
		))
		return
	}
	err := h.service.SetUserActive(gymID, id, activeReq.Active)
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
	if activeReq.Active {
		statusMsg = "activated"
	}
	response.WriteAPISuccess(w, fmt.Sprintf("User %s successfully", statusMsg), nil)
}

// Platform Admin Methods - Allow specifying gym context

func (h *UsersHandler) GetUsersByGymID(w http.ResponseWriter, r *http.Request) {
	// Security check: Only platform admins can access users from any gym
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can access users from any gym",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID is required",
			nil,
		))
		return
	}

	users, err := h.service.GetAllUsers(gymID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
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

func (h *UsersHandler) RegisterUserInGym(w http.ResponseWriter, r *http.Request) {
	// Security check: Only platform admins can create users in any gym
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can create users in any gym",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID is required",
			nil,
		))
		return
	}

	var userDTO dto.UserCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	// Validate required fields
	if userDTO.Username == "" || userDTO.Email == "" || userDTO.Password == "" || userDTO.Role == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Missing required user fields",
			nil,
		))
		return
	}

	// Platform admin can create any type of user
	userID, err := h.service.RegisterUser(gymID, &userDTO)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
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
	response.WriteAPISuccess(w, "User created successfully", map[string]string{"id": *userID})
}

func (h *UsersHandler) GetUserByIDInGym(w http.ResponseWriter, r *http.Request) {
	// Security check: Only platform admins can access users from any gym
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can access users from any gym",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	userID := chi.URLParam(r, "id")

	if gymID == "" || userID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID and User ID are required",
			nil,
		))
		return
	}

	user, err := h.service.GetUserByID(gymID, userID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
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
	response.WriteAPISuccess(w, "User retrieved successfully", user)
}

func (h *UsersHandler) UpdateUserInGym(w http.ResponseWriter, r *http.Request) {
	// Security check: Only platform admins can update users in any gym
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can update users in any gym",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	userID := chi.URLParam(r, "id")

	if gymID == "" || userID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID and User ID are required",
			nil,
		))
		return
	}

	var userDTO dto.UserUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	err := h.service.UpdateUser(gymID, userID, &userDTO)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
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

func (h *UsersHandler) DeleteUserInGym(w http.ResponseWriter, r *http.Request) {
	// Security check: Only platform admins can delete users from any gym
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can delete users from any gym",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	userID := chi.URLParam(r, "id")

	if gymID == "" || userID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID and User ID are required",
			nil,
		))
		return
	}

	err := h.service.DeleteUser(gymID, userID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
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
	response.WriteAPISuccess(w, "User deleted successfully", nil)
}

func (h *UsersHandler) VerifyUserInGym(w http.ResponseWriter, r *http.Request) {
	// Security check: Only platform admins can verify users in any gym
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can verify users in any gym",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	userID := chi.URLParam(r, "id")

	if gymID == "" || userID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID and User ID are required",
			nil,
		))
		return
	}

	err := h.service.VerifyUser(gymID, userID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
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

func (h *UsersHandler) SetUserActiveInGym(w http.ResponseWriter, r *http.Request) {
	// Security check: Only platform admins can set user status in any gym
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Access denied: Only platform administrators can set user status in any gym",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	userID := chi.URLParam(r, "id")

	if gymID == "" || userID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID and User ID are required",
			nil,
		))
		return
	}

	var activeReq struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&activeReq); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	err := h.service.SetUserActive(gymID, userID, activeReq.Active)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
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
	if activeReq.Active {
		statusMsg = "activated"
	}
	response.WriteAPISuccess(w, fmt.Sprintf("User %s successfully", statusMsg), nil)
}
