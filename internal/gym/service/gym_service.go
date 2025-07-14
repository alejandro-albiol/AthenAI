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
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to check domain existence", err)
	}
	if err == nil {
		return "", apierror.New(errorcode_enum.CodeConflict, "Gym with this domain already exists", nil)
	}

	gymID, err := s.repository.CreateGym(createDTO)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create gym", err)
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to connect to database", err)
	}
	defer db.Close()

	// Use the domain name for the schema, not the gym ID
	err = database.CreateTenantSchema(db, createDTO.Domain)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create tenant schema", err)
	}

	return gymID, nil
}

func (s *GymService) GetGymByID(id string) (dto.GymResponseDTO, error) {
	gym, err := s.repository.GetGymByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found", nil)
		}
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get gym", err)
	}

	return gym, nil
}

func (s *GymService) GetGymByDomain(domain string) (dto.GymResponseDTO, error) {
	gym, err := s.repository.GetGymByDomain(domain)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found", nil)
		}
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get gym", err)
	}

	return gym, nil
}

func (s *GymService) GetAllGyms() ([]dto.GymResponseDTO, error) {
	gyms, err := s.repository.GetAllGyms()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get gyms", err)
	}

	if len(gyms) == 0 {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "No gyms found", nil)
	}

	return gyms, nil
}

func (s *GymService) UpdateGym(id string, updateDTO dto.GymUpdateDTO) (dto.GymResponseDTO, error) {

	_, err := s.repository.GetGymByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found", nil)
		}
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to check gym existence", err)
	}

	updatedGym, err := s.repository.UpdateGym(id, updateDTO)
	if err != nil {
		return dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update gym", err)
	}

	return updatedGym, nil
}

func (s *GymService) SetGymActive(id string, active bool) error {

	_, err := s.repository.GetGymByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierror.New(errorcode_enum.CodeNotFound, "Gym not found", nil)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check gym existence", err)
	}

	err = s.repository.SetGymActive(id, active)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update gym status", err)
	}

	return nil
}

func (s *GymService) DeleteGym(id string) error {
	err := s.repository.DeleteGym(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierror.New(errorcode_enum.CodeNotFound, "Gym not found or already deleted", nil)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete gym", err)
	}

	return nil
}
