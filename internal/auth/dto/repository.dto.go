package dto

import "time"

// Repository DTOs for internal database operations

// AdminAuthDTO - Platform admin authentication data from public.admin table
type AdminAuthDTO struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	Email        string     `json:"email"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TenantUserAuthDTO - Tenant user authentication data from tenant schema
type TenantUserAuthDTO struct {
	ID                 string     `json:"id"`
	Username           string     `json:"username"`
	PasswordHash       string     `json:"password_hash"`
	Email              string     `json:"email"`
	Role               string     `json:"role"`
	VerificationStatus string     `json:"verification_status"`
	IsActive           bool       `json:"is_active"`
	GymID              string     `json:"gym_id"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

// GymDomainDTO - Gym domain validation data from gym table
type GymDomainDTO struct {
	ID       string `json:"id"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Schema   string `json:"schema"`
	IsActive bool   `json:"is_active"`
}

// RefreshTokenDTO - Stored refresh token data from refresh_tokens table
type RefreshTokenDTO struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	UserType  UserType  `json:"user_type"`
	GymID     *string   `json:"gym_id,omitempty"` // For tenant users
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
