package interfaces

import (
	"time"

	"github.com/alejandro-albiol/athenai/internal/user/enum"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     userrole_enum.UserRole `json:"role,omitempty"`
	GymID    string `json:"gym_id"`
	Verified bool   `json:"verified,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
