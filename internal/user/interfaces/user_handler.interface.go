package interfaces

import "github.com/alejandro-albiol/athenai/internal/user/dto"

type UserHandler interface {
	// RegisterUser handles user registration.
	RegisterUser(dto dto.UserCreationDTO) (User, error)
	// GetUserByID handles retrieving a user by ID.
	GetUserByID(id string) (User, error)
	// GetUserByUsername handles retrieving a user by username.
	GetUserByUsername(username string) (User, error)
	// GetUserByEmail handles retrieving a user by email.
	GetUserByEmail(email string) (User, error)
	// GetAllUsers retrieves all users.
	GetAllUsers() ([]User, error)
	// GetPasswordHashByUsername retrieves the password hash for a given username.
	GetPasswordHashByUsername(username string) (string, error)
	// UpdateUser handles updating an existing user.
	UpdateUser(id string, dto dto.UserUpdateDTO) (User, error)
	// UpdatePassword handles updating the password for a user.
	UpdatePassword(id string, newPasswordHash string) error
	// DeleteUser handles removing a user by ID.
	DeleteUser(id string) error
}
