package interfaces

import dto "github.com/alejandro-albiol/athenai/internal/user/dto"

type UserService interface {
	// RegisterUser registers a new user.
	RegisterUser(user dto.UserCreationDTO) error
	// GetUserByID retrieves a user by their ID.
	GetUserByID(id string) (dto.UserResponseDTO, error)
	// GetUserByUsername retrieves a user by username.
	GetUserByUsername(username string) (dto.UserResponseDTO, error)
	// GetUserByEmail retrieves a user by email.
	GetUserByEmail(email string) (dto.UserResponseDTO, error)
	// GetAllUsers retrieves all users.
	GetAllUsers() ([]dto.UserResponseDTO, error)
	// GetPasswordHashByUsername retrieves the password hash for a given username.
	GetPasswordHashByUsername(username string) (string, error)
	// UpdateUser updates an existing user.
	UpdateUser(user dto.UserUpdateDTO) error
	// UpdatePassword updates the password for a user.
	UpdatePassword(id string, newPasswordHash string) error
	// DeleteUser removes a user by ID.
	DeleteUser(id string) error
}