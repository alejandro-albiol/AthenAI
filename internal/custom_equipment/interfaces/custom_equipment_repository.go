package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"

// CustomEquipmentRepository defines DB operations for custom equipment
// All methods require gymID for multi-tenancy

type CustomEquipmentRepository interface {
	Create(gymID string, equipment *dto.CreateCustomEquipmentDTO) error
	GetByID(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error)
	GetByName(gymID, name string) (*dto.ResponseCustomEquipmentDTO, error)
	List(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error)
	Update(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error
	Delete(gymID, id string) error
}
