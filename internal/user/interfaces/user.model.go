package interfaces

import (
	"time"

	"github.com/alejandro-albiol/athenai/internal/user/enum"
)

type User struct {
	ID               string        `json:"id"`
	Username         string        `json:"username"`
	Email            string        `json:"email"`
	Password         string        `json:"password,omitempty"`
	Role             enum.UserRole `json:"role,omitempty"`
	Verified         bool          `json:"verified,omitempty"`
	Description      *string       `json:"description,omitempty"`
	TrainingPhase    string        `json:"training_phase"`
	Motivation       string        `json:"motivation"`
	SpecialSituation string        `json:"special_situation"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}
