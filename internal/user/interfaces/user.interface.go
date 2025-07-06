package interfaces

import "github.com/alejandro-albiol/athenai/internal/user/enum"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     userrole_enum.UserRole `json:"role"`
}
