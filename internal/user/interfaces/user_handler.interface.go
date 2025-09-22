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

	// Platform admin methods - allow specifying gym context
	// GetUsersByGymID retrieves all users for a specific gym (platform admin only).
	GetUsersByGymID(w http.ResponseWriter, r *http.Request)
	// RegisterUserInGym handles user registration in a specific gym (platform admin only).
	RegisterUserInGym(w http.ResponseWriter, r *http.Request)
	// GetUserByIDInGym handles retrieving a user by ID in a specific gym (platform admin only).
	GetUserByIDInGym(w http.ResponseWriter, r *http.Request)
	// UpdateUserInGym handles updating an existing user in a specific gym (platform admin only).
	UpdateUserInGym(w http.ResponseWriter, r *http.Request)
	// DeleteUserInGym handles removing a user by ID in a specific gym (platform admin only).
	DeleteUserInGym(w http.ResponseWriter, r *http.Request)
	// VerifyUserInGym marks a user as verified in a specific gym (platform admin only).
	VerifyUserInGym(w http.ResponseWriter, r *http.Request)
	// SetUserActiveInGym sets a user's active status in a specific gym (platform admin only).
	SetUserActiveInGym(w http.ResponseWriter, r *http.Request)
}
