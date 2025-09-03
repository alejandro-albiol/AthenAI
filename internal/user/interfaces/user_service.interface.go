package interfaces

import dto "github.com/alejandro-albiol/athenai/internal/user/dto"

type UserService interface {
	// RegisterUser registers a new user.
	RegisterUser(gymID string, user *dto.UserCreationDTO) (string, error)
	// GetUserByID retrieves a user by their ID.
	GetUserByID(gymID, id string) (*dto.UserResponseDTO, error)
	// GetUserByUsername retrieves a user by username.
	GetUserByUsername(gymID, username string) (*dto.UserResponseDTO, error)
	// GetUserByEmail retrieves a user by email.
	GetUserByEmail(gymID, email string) (*dto.UserResponseDTO, error)
	// GetAllUsers retrieves all users for a gym.
	GetAllUsers(gymID string) ([]*dto.UserResponseDTO, error)
	// GetPasswordHashByUsername retrieves the password hash for a given username.
	GetPasswordHashByUsername(gymID, username string) (string, error)
	// UpdateUser updates an existing user.
	UpdateUser(gymID string, id string, user *dto.UserUpdateDTO) error
	// UpdatePassword updates the password for a user.
	UpdatePassword(gymID, id string, newPasswordHash string) error
	// DeleteUser removes a user by ID.
	DeleteUser(gymID, id string) error
	// VerifyUser marks a user as verified.
	VerifyUser(gymID, userID string) error
	// SetUserActive sets a user's active status.
	SetUserActive(gymID, userID string, active bool) error
}
