package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	authdto "github.com/alejandro-albiol/athenai/internal/auth/dto"
	authinterfaces "github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type AuthHandler struct {
	authService authinterfaces.AuthServiceInterface
}

func NewAuthHandler(authService authinterfaces.AuthServiceInterface) authinterfaces.AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq *authdto.LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request body",
			err,
		))
		return
	}

	loginResp, apiErr := h.authService.Login(r, loginReq)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Login successful", loginResp)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Authorization header missing",
			nil,
		))
		return
	}

	// Improved token extraction with trimming
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		// TrimPrefix didn't change the string, meaning "Bearer " prefix was not found
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid authorization header format",
			nil,
		))
		return
	}

	token = strings.TrimSpace(token)
	if token == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Empty token provided",
			nil,
		))
		return
	}

	validationResp, apiErr := h.authService.ValidateToken(token)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Token validation successful", validationResp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshReq *authdto.RefreshTokenRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request body",
			err,
		))
		return
	}

	loginResp, apiErr := h.authService.RefreshToken(refreshReq)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Token refreshed successfully", loginResp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var logoutReq *authdto.LogoutRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request body",
			err,
		))
		return
	}

	apiErr := h.authService.Logout(logoutReq)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Logout successful", nil)
}
