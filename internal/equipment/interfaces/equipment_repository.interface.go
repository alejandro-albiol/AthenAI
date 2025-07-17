package interfaces

import "github.com/alejandro-albiol/athenai/internal/equipment/dto"

type EquipmentRepository interface {
	CreateEquipment(equipment dto.EquipmentCreationDTO) (string, error)
	GetEquipmentByID(id string) (dto.EquipmentResponseDTO, error)
	GetAllEquipment() ([]dto.EquipmentResponseDTO, error)
	UpdateEquipment(id string, update dto.EquipmentUpdateDTO) (dto.EquipmentResponseDTO, error)
	DeleteEquipment(id string) error
}
