package dto

import (
	"time"

	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
)

type UserResponseDTO struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	Role      userrole_enum.UserRole `json:"role"`
	GymID     string                 `json:"gym_id"`
	Verified  bool                   `json:"verified"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
