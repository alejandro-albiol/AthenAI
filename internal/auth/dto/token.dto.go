package dto

import "github.com/golang-jwt/jwt/v5"

// TokenValidationResponseDTO - Response from token validation
type TokenValidationResponseDTO struct {
	Valid   bool      `json:"valid"`
	Claims  ClaimsDTO `json:"claims,omitempty"`
	Message string    `json:"message"`
}

// ClaimsDTO - JWT token claims (implements jwt.Claims)
type ClaimsDTO struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	UserType string  `json:"user_type"` // "platform_admin" or "tenant_user"
	GymID    *string `json:"gym_id,omitempty"`
	Role     *string `json:"role,omitempty"` // admin, user, guest
	IsActive bool    `json:"is_active"`

	jwt.RegisteredClaims
}

// RefreshTokenRequestDTO - Request to refresh access token
type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequestDTO - Request to logout and revoke tokens
type LogoutRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
