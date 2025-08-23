package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/interfaces"
)

// CustomEquipmentServiceImpl implements interfaces.CustomEquipmentService

type CustomEquipmentServiceImpl struct {
	Repo interfaces.CustomEquipmentRepository
}

func (s *CustomEquipmentServiceImpl) CreateCustomEquipment(gymID string, equipment *dto.CreateCustomEquipmentDTO) error {
	return s.Repo.Create(gymID, equipment)
}

func (s *CustomEquipmentServiceImpl) GetCustomEquipmentByID(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error) {
	return s.Repo.GetByID(gymID, id)
}

func (s *CustomEquipmentServiceImpl) ListCustomEquipment(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
	return s.Repo.List(gymID)
}

func (s *CustomEquipmentServiceImpl) UpdateCustomEquipment(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error {
	return s.Repo.Update(gymID, equipment)
}

func (s *CustomEquipmentServiceImpl) DeleteCustomEquipment(gymID, id string) error {
	return s.Repo.Delete(gymID, id)
}
