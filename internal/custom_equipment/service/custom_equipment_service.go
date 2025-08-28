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

func (s *CustomEquipmentService) CreateCustomEquipment(gymID string, equipment *dto.CreateCustomEquipmentDTO) error {
	existingEquipment, err := s.Repo.GetByID(gymID, equipment.Name)
	if err != nil && err != sql.ErrNoRows {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check existing equipment", err)
	}
	if existingEquipment != nil {
		return apierror.New(errorcode_enum.CodeConflict, "Equipment already exists", nil)
	}
	err = s.Repo.Create(gymID, equipment)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create equipment", err)
	}
	return nil
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

func (s *CustomEquipmentService) ListCustomEquipment(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
	equipment, err := s.Repo.List(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list equipment", err)
	}
	return equipment, nil
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
