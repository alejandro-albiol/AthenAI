package dto

import userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"

type UserResponseDTO struct {
	ID       string                 `json:"id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Role     userrole_enum.UserRole `json:"role"`
	GymID    string                 `json:"gym_id"`
}
