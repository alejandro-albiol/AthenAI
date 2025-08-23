package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"

// CustomEquipmentService defines business logic for custom equipment

type CustomEquipmentService interface {
	CreateCustomEquipment(gymID string, equipment *dto.CreateCustomEquipmentDTO) error
	GetCustomEquipmentByID(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error)
	ListCustomEquipment(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error)
	UpdateCustomEquipment(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error
	DeleteCustomEquipment(gymID, id string) error
}
