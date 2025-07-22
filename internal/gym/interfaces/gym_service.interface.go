package interfaces

import "github.com/alejandro-albiol/athenai/internal/gym/dto"

// GymService defines the interface for gym-related business logic.
// It handles domain rules, error mapping from repository layer,
// and coordinates complex operations involving multiple repositories.
type GymService interface {
	// CreateGym validates and creates a new gym.
	// Returns the created gym's ID on success, or a domain error if validation fails
	// or there are conflicts (e.g., duplicate domain).
	CreateGym(gym dto.GymCreationDTO) (string, error)

	// GetGymByID retrieves a gym by its unique identifier.
	// Maps database errors to domain errors and validates access permissions.
	GetGymByID(id string) (dto.GymResponseDTO, error)

	// GetGymByName retrieves a gym by its unique name.
	// Maps database errors to domain errors and validates access permissions.
	GetGymByName(name string) (dto.GymResponseDTO, error)

	// GetAllGyms retrieves all active gyms.
	// Maps database errors to domain errors and applies any business filters.
	GetAllGyms() ([]dto.GymResponseDTO, error)

	// UpdateGym validates and updates an existing gym.
	// Returns a domain error if validation fails, the gym doesn't exist,
	// or there are conflicts (e.g., duplicate domain).
	UpdateGym(id string, gym dto.GymUpdateDTO) (dto.GymResponseDTO, error)

	// SetGymActive updates the active status of a gym.
	// Returns a domain error if the gym doesn't exist or the operation fails.
	SetGymActive(id string, active bool) error

	// DeleteGym soft deletes a gym and handles any cascade operations.
	// Returns a domain error if the gym doesn't exist or the operation fails.
	DeleteGym(id string) error
}
