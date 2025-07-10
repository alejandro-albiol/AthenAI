package interfaces

import "github.com/alejandro-albiol/athenai/internal/gym/dto"

type GymHandler interface {
	CreateGym(gym dto.GymCreationDTO) (string, error)
	GetGym(id string) (dto.GymResponseDTO, error)
	GetGymByDomain(domain string) (dto.GymResponseDTO, error)
	GetAllGyms() ([]dto.GymResponseDTO, error)
	UpdateGym(id string, gym dto.GymResponseDTO) error
	SetGymActive(id string, active bool) error
	DeleteGym(id string) error
}
