package dto

// LoginRequestDTO - Single login DTO that works for both admin and tenant users
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponseDTO - Authentication response
type LoginResponseDTO struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	UserInfo     UserInfoDTO `json:"user_info"`
}

// UserInfoDTO - User information returned after login
type UserInfoDTO struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	UserType string  `json:"user_type"`        // "platform_admin" or "tenant_user"
	Role     *string `json:"role,omitempty"`   // For tenant users: admin, user, guest
	GymID    *string `json:"gym_id,omitempty"` // For tenant users
}
