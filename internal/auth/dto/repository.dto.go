package dto

import "time"

// Repository DTOs for internal database operations

// AdminAuthDTO - Platform admin authentication data from public.admin table
type AdminAuthDTO struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// TenantUserAuthDTO - Tenant user authentication data from {domain}.users table
type TenantUserAuthDTO struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`        // admin, user, guest
	IsVerified bool      `json:"is_verified"` // verification status
	IsActive   bool      `json:"is_active"`
	GymID      string    `json:"gym_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// RefreshTokenDTO - Stored refresh token data
type RefreshTokenDTO struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	UserType  string    `json:"user_type"`        // "platform_admin" or "tenant_user"
	GymID     *string   `json:"gym_id,omitempty"` // For tenant users
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
