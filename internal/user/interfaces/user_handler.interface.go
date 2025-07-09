package interfaces

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
)

type UserHandler interface {
	// RegisterUser handles user registration.
	RegisterUser(w http.ResponseWriter, r *http.Request, gymID string)
	// GetUserByID handles retrieving a user by ID.
	GetUserByID(w http.ResponseWriter, gymID, id string)
	// GetUserByUsername handles retrieving a user by username.
	GetUserByUsername(w http.ResponseWriter, gymID, username string)
	// GetUserByEmail handles retrieving a user by email.
	GetUserByEmail(w http.ResponseWriter, gymID, email string)
	// GetAllUsers retrieves all users.
	GetAllUsers(w http.ResponseWriter, gymID string)
	// UpdateUser handles updating an existing user.
	UpdateUser(w http.ResponseWriter, gymID, id string, user dto.UserUpdateDTO)
	// DeleteUser handles removing a user by ID.
	DeleteUser(w http.ResponseWriter, gymID, id string)
	// VerifyUser marks a user as verified.
	VerifyUser(w http.ResponseWriter, gymID, id string)
	// SetUserActive sets a user's active status.
	SetUserActive(w http.ResponseWriter, gymID, id string, active bool)
}
