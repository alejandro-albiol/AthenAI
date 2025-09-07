package dto

import (
	"time"

	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
)

type UserResponseDTO struct {
	ID               string                 `json:"id"`
	Username         string                 `json:"username"`
	Email            string                 `json:"email"`
	Role             userrole_enum.UserRole `json:"role"`
	Verified         bool                   `json:"verified"`
	IsActive         bool                   `json:"is_active"`
	Phone            string                 `json:"phone"`
	Description      *string                `json:"description,omitempty"`
	TrainingPhase    string                 `json:"training_phase"`
	Motivation       string                 `json:"motivation"`
	SpecialSituation string                 `json:"special_situation"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}
