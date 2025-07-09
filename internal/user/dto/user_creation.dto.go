package dto

import userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"

type UserCreationDTO struct {
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Role     userrole_enum.UserRole `json:"role,omitempty"`
}
