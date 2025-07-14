package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/alejandro-albiol/athenai/internal/auth/dto"
	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type AuthHandler struct {
	service interfaces.AuthServiceInterface
}

func NewAuthHandler(service interfaces.AuthServiceInterface) interfaces.AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login handles both platform admin and tenant user authentication based on subdomain
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		apiErr := apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		)
		response.WriteAPIError(w, apiErr)
		return
	}

	// Get gym ID from middleware (X-Gym-ID header)
	gymID := middleware.GetGymID(r)

	var loginResponse dto.LoginResponseDTO
	var err error

	if gymID == "" {
		// Platform admin login (no gym ID provided)
		loginResponse, err = h.service.LoginAdmin(credentials)
	} else {
		// Tenant user login (gym ID provided via X-Gym-ID header)
		loginResponse, err = h.service.LoginTenantUser(gymID, credentials)
	}

	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		// Fallback for non-APIErrors
		internalErr := apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error during login",
			err,
		)
		response.WriteAPIError(w, internalErr)
		return
	}

	response.WriteAPISuccess(w, "Login successful", loginResponse)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshReq dto.RefreshTokenRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		apiErr := apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		)
		response.WriteAPIError(w, apiErr)
		return
	}

	loginResponse, err := h.service.RefreshToken(refreshReq)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		// Fallback for non-APIErrors
		internalErr := apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error during token refresh",
			err,
		)
		response.WriteAPIError(w, internalErr)
		return
	}

	response.WriteAPISuccess(w, "Token refreshed successfully", loginResponse)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var logoutReq dto.LogoutRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
		apiErr := apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		)
		response.WriteAPIError(w, apiErr)
		return
	}

	logoutResponse, err := h.service.Logout(logoutReq)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		// Fallback for non-APIErrors
		internalErr := apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error during logout",
			err,
		)
		response.WriteAPIError(w, internalErr)
		return
	}

	response.WriteAPISuccess(w, "Logout successful", logoutResponse)
}

// ValidateToken handles token validation
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		apiErr := apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Authorization header required",
			nil,
		)
		response.WriteAPIError(w, apiErr)
		return
	}

	// Remove "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		apiErr := apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid authorization format",
			nil,
		)
		response.WriteAPIError(w, apiErr)
		return
	}

	validationResponse, err := h.service.ValidateToken(token)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			response.WriteAPIError(w, apiErr)
			return
		}
		// Fallback for non-APIErrors
		internalErr := apierror.New(
			errorcode_enum.CodeInternal,
			"Internal server error during token validation",
			err,
		)
		response.WriteAPIError(w, internalErr)
		return
	}

	response.WriteAPISuccess(w, "Token is valid", validationResponse)
}
