package interfaces

import "github.com/alejandro-albiol/athenai/internal/gym/dto"

// GymRepository defines the interface for gym data persistence operations.
// It handles all database operations and returns raw database errors
// which will be mapped to domain errors by the service layer.
type GymRepository interface {
	// CreateGym persists a new gym to the database.
	// Returns raw database errors without any domain error mapping.
	CreateGym(gym dto.GymCreationDTO) (string, error)

	// GetGymByID retrieves a gym from the database by its ID.
	// Returns sql.ErrNoRows if not found, or other raw database errors.
	GetGymByID(id string) (dto.GymResponseDTO, error)

	// GetGymByDomain retrieves a gym from the database by its domain.
	// Returns sql.ErrNoRows if not found, or other raw database errors.
	GetGymByDomain(domain string) (dto.GymResponseDTO, error)

	// GetAllGyms retrieves all active gyms from the database.
	// Returns raw database errors without any domain error mapping.
	GetAllGyms() ([]dto.GymResponseDTO, error)

	// UpdateGym updates an existing gym in the database.
	// Returns sql.ErrNoRows if the gym doesn't exist, or other raw database errors.
	UpdateGym(id string, gym dto.GymUpdateDTO) (dto.GymResponseDTO, error)

	// SetGymActive updates the active status of a gym.
	// Returns sql.ErrNoRows if the gym doesn't exist, or other raw database errors.
	SetGymActive(id string, active bool) error

	// DeleteGym soft deletes a gym by setting its deleted_at timestamp.
	// Returns sql.ErrNoRows if the gym doesn't exist, or other raw database errors.
	DeleteGym(id string) error
}
