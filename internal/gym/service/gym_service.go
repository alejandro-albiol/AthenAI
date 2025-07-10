package service

import (
	"database/sql"
	"errors"

	"github.com/alejandro-albiol/athenai/internal/database"
	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type GymService struct {
	repository interfaces.GymRepository
}

func NewGymService(repository interfaces.GymRepository) *GymService {
	return &GymService{repository: repository}
}

func (s *GymService) CreateGym(createDTO dto.GymCreationDTO) (string, error) {
	_, err := s.repository.GetGymByDomain(createDTO.Domain)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to check domain existence")
	}
	if err == nil {
		return "", apierror.New(errorcode_enum.CodeConflict, "Gym with this domain already exists")
	}

	domain, err := s.repository.CreateGym(createDTO)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create gym")
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to connect to database")
	}
	defer db.Close()
	
	database.CreateTenantSchema(db, domain)

	return domain, nil
}

func (s *GymService) GetGymByID(id string) (dto.GymResponseDTO, error) {
	gym, err := s.repository.GetGymByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found")
		}
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get gym")
	}

	return gym, nil
}

func (s *GymService) GetGymByDomain(domain string) (dto.GymResponseDTO, error) {
	gym, err := s.repository.GetGymByDomain(domain)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found")
		}
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get gym")
	}

	return gym, nil
}

func (s *GymService) GetAllGyms() ([]dto.GymResponseDTO, error) {
	gyms, err := s.repository.GetAllGyms()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get gyms")
	}

	return gyms, nil
}

func (s *GymService) UpdateGym(id string, updateDTO dto.GymUpdateDTO) (dto.GymResponseDTO, error) {
	// Check if gym exists
	_, err := s.repository.GetGymByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found")
		}
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to check gym existence")
	}

	updatedGym, err := s.repository.UpdateGym(id, updateDTO)
	if err != nil {
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update gym")
	}

	return updatedGym, nil
}

func (s *GymService) SetGymActive(id string, active bool) error {

	_, err := s.repository.GetGymByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierror.New(errorcode_enum.CodeNotFound, "Gym not found")
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check gym existence")
	}

	err = s.repository.SetGymActive(id, active)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update gym status")
	}

	return nil
}

func (s *GymService) DeleteGym(id string) error {
	// Check if gym exists
	_, err := s.repository.GetGymByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierror.New(errorcode_enum.CodeNotFound, "Gym not found")
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check gym existence")
	}

	err = s.repository.DeleteGym(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete gym")
	}

	return nil
}
