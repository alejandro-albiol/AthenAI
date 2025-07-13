package dto

// TokenValidationResponseDTO - Response from token validation
type TokenValidationResponseDTO struct {
	Valid   bool       `json:"valid"`
	Claims  *ClaimsDTO `json:"claims,omitempty"`
	Message string     `json:"message"`
}

// ClaimsDTO - JWT token claims
type ClaimsDTO struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	UserType UserType `json:"user_type"`

	// Platform admin claims
	IsActive *bool `json:"is_active,omitempty"`

	// Tenant user claims
	GymDomain          *string `json:"gym_domain,omitempty"`
	Role               *string `json:"role,omitempty"`
	VerificationStatus *string `json:"verification_status,omitempty"`

	ExpiresAt int64 `json:"exp"`
	IssuedAt  int64 `json:"iat"`
}

// RefreshTokenRequestDTO - Request to refresh access token
type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequestDTO - Request to logout and revoke tokens
type LogoutRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutResponseDTO - Response after logout
type LogoutResponseDTO struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
