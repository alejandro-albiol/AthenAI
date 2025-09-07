package service

import (
	"github.com/alejandro-albiol/athenai/internal/equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type EquipmentService struct {
	repo interfaces.EquipmentRepository
}

func NewEquipmentService(repo interfaces.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: repo}
}

func (s *EquipmentService) CreateEquipment(equipment *dto.EquipmentCreationDTO) (*string, error) {
	// Uniqueness check by name (assuming name is unique)
	allEquipment, err := s.repo.GetAllEquipment()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to check equipment name uniqueness", err)
	}
	for _, eq := range allEquipment {
		if eq.Name == equipment.Name {
			return nil, apierror.New(errorcode_enum.CodeConflict, "Equipment with this name already exists", nil)
		}
	}
	id, err := s.repo.CreateEquipment(equipment)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create equipment", err)
	}
	return id, nil
}

func (s *EquipmentService) GetEquipmentByID(id string) (*dto.EquipmentResponseDTO, error) {
	equipment, err := s.repo.GetEquipmentByID(id)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "Equipment with ID not found", err)
	}
	return equipment, nil
}

func (s *EquipmentService) GetAllEquipment() ([]*dto.EquipmentResponseDTO, error) {
	equipmentList, err := s.repo.GetAllEquipment()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve equipment", err)
	}
	return equipmentList, nil
}

func (s *EquipmentService) UpdateEquipment(id string, update *dto.EquipmentUpdateDTO) (*dto.EquipmentResponseDTO, error) {
	// Check existence first
	_, err := s.repo.GetEquipmentByID(id)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "Equipment with ID not found", err)
	}
	equipment, err := s.repo.UpdateEquipment(id, update)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to update equipment", err)
	}
	return equipment, nil
}

func (s *EquipmentService) DeleteEquipment(id string) error {
	// Check existence first
	_, err := s.repo.GetEquipmentByID(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeNotFound, "Equipment with ID not found", err)
	}
	err = s.repo.DeleteEquipment(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete equipment", err)
	}
	return nil
}
