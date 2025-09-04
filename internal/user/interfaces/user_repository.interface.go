package interfaces

import dto "github.com/alejandro-albiol/athenai/internal/user/dto"

type UserRepository interface {
	// CreateUser creates a new user in the database for a specific gym.
	CreateUser(gymID string, user *dto.UserCreationDTO) (*string, error)
	// GetUserByID retrieves a user by their ID and gym.
	GetUserByID(gymID, id string) (*dto.UserResponseDTO, error)
	// GetUserByUsername retrieves a user by their username and gym.
	GetUserByUsername(gymID, username string) (*dto.UserResponseDTO, error)
	// GetUserByEmail retrieves a user by their email and gym.
	GetUserByEmail(gymID, email string) (*dto.UserResponseDTO, error)
	// GetAllUsers retrieves all users for a specific gym.
	GetAllUsers(gymID string) ([]*dto.UserResponseDTO, error)
	// GetPasswordHashByUsername retrieves the password hash for a given username and gym.
	GetPasswordHashByUsername(gymID, username string) (string, error)
	// UpdateUser updates an existing user in the database for a specific gym.
	UpdateUser(gymID, id string, user *dto.UserUpdateDTO) error
	// UpdatePassword updates the password for a user in a specific gym.
	UpdatePassword(gymID, id string, newPasswordHash string) error
	// DeleteUser removes a user from the database for a specific gym.
	DeleteUser(gymID, id string) error
	// VerifyUser marks a user as verified in the database for a specific gym.
	VerifyUser(gymID, userID string) error
	// SetUserActive sets a user's active status in the database for a specific gym.
	SetUserActive(gymID, userID string, active bool) error
}
