package interfaces

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/auth/dto"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
)

// AuthServiceInterface defines authentication business logic with simplified approach
type AuthServiceInterface interface {
	// Single login method that checks X-Gym-ID header to determine authentication type
	Login(r *http.Request, loginReq dto.LoginRequestDTO) (*dto.LoginResponseDTO, *apierror.APIError)

	// Token operations
	ValidateToken(token string) (*dto.TokenValidationResponseDTO, *apierror.APIError)
	RefreshToken(refreshReq dto.RefreshTokenRequestDTO) (*dto.LoginResponseDTO, *apierror.APIError)
	Logout(logoutReq dto.LogoutRequestDTO) *apierror.APIError

}
