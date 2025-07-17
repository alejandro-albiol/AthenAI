package interfaces

import "github.com/alejandro-albiol/athenai/internal/equipment/dto"

// EquipmentService defines the contract for business logic on equipment

type EquipmentService interface {
	CreateEquipment(equipment dto.EquipmentCreationDTO) (string, error)
	GetEquipmentByID(id string) (dto.EquipmentResponseDTO, error)
	GetAllEquipment() ([]dto.EquipmentResponseDTO, error)
	UpdateEquipment(id string, update dto.EquipmentUpdateDTO) (dto.EquipmentResponseDTO, error)
	DeleteEquipment(id string) error
}
