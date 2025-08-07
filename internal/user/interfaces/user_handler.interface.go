package interfaces

import (
	"net/http"
)

type UserHandler interface {
	// RegisterUser handles user registration.
	RegisterUser(w http.ResponseWriter, r *http.Request)
	// GetUserByID handles retrieving a user by ID.
	GetUserByID(w http.ResponseWriter, r *http.Request)
	// GetUserByUsername handles retrieving a user by username.
	GetUserByUsername(w http.ResponseWriter, r *http.Request)
	// GetUserByEmail handles retrieving a user by email.
	GetUserByEmail(w http.ResponseWriter, r *http.Request)
	// GetAllUsers retrieves all users.
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	// UpdateUser handles updating an existing user.
	UpdateUser(w http.ResponseWriter, r *http.Request)
	// DeleteUser handles removing a user by ID.
	DeleteUser(w http.ResponseWriter, r *http.Request)
	// VerifyUser marks a user as verified.
	VerifyUser(w http.ResponseWriter, r *http.Request)
	// SetUserActive sets a user's active status.
	SetUserActive(w http.ResponseWriter, r *http.Request)
}
