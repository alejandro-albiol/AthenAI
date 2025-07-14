package dto

import "time"

// UserType enum for JWT claims
type UserType string

const (
	UserTypePlatformAdmin UserType = "platform_admin"
	UserTypeTenantUser    UserType = "tenant_user"
)

// LoginRequestDTO - Single login DTO for both admin and tenant users
type LoginRequestDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponseDTO - Authentication response
type LoginResponseDTO struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    time.Time   `json:"expires_at"`
	UserInfo     UserInfoDTO `json:"user_info"`
}

// UserInfoDTO - User information returned after login
type UserInfoDTO struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"`

	// For platform admins
	IsActive *bool `json:"is_active,omitempty"`

	// For tenant users
	GymID              *string `json:"gym_id,omitempty"`
	Role               *string `json:"role,omitempty"`
	VerificationStatus *string `json:"verification_status,omitempty"`
}
