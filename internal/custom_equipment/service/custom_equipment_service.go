package service

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomEquipmentService struct {
	Repo interfaces.CustomEquipmentRepository
}

func NewCustomEquipmentService(repo interfaces.CustomEquipmentRepository) *CustomEquipmentService {
	return &CustomEquipmentService{Repo: repo}
}

func (s *CustomEquipmentService) CreateCustomEquipment(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error) {
	if equipment.Name == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "Name is required", nil)
	}
	if equipment.CreatedBy == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "CreatedBy is required", nil)
	}
	existingEquipment, err := s.Repo.GetByName(gymID, equipment.Name)
	if err != nil && err != sql.ErrNoRows {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to check existing equipment", err)
	}
	if existingEquipment != nil && existingEquipment.ID != "" {
		return nil, apierror.New(errorcode_enum.CodeConflict, "Equipment already exists", nil)
	}
	id, err := s.Repo.Create(gymID, equipment)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create equipment", err)
	}
	return id, nil
}

func (s *CustomEquipmentService) GetCustomEquipmentByID(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error) {
	equipment, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get equipment", err)
	}
	return equipment, nil
}

func (s *CustomEquipmentService) GetCustomEquipmentByName(gymID, name string) (*dto.ResponseCustomEquipmentDTO, error) {
	equipment, err := s.Repo.GetByName(gymID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get equipment", err)
	}
	return equipment, nil
}

func (s *CustomEquipmentService) ListCustomEquipment(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
	equipment, err := s.Repo.List(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list equipment", err)
	}
	var result []*dto.ResponseCustomEquipmentDTO
	result = append(result, equipment...)
	return result, nil
}

func (s *CustomEquipmentService) UpdateCustomEquipment(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error {
	err := s.Repo.Update(gymID, equipment)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update equipment", err)
	}
	return nil
}

func (s *CustomEquipmentService) DeleteCustomEquipment(gymID, id string) error {
	err := s.Repo.Delete(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete equipment", err)
	}
	return nil
}
