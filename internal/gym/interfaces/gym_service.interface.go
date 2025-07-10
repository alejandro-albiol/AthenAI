package interfaces

import (
	"github.com/alejandro-albiol/athenai/internal/gym/dto"
)

type GymService interface {
	CreateGym(gym dto.GymCreationDTO) (Gym, error)
	GetGymByID(id string) (Gym, error)
	GetGymByDomain(domain string) (Gym, error)
	GetAllGyms() ([]Gym, error)
	UpdateGym(gym dto.GymUpdateDTO) (Gym, error)
	SetGymActive(id string, active bool) error
	DeleteGym(id string) error
}
