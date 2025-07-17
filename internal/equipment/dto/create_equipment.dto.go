package dto

// EquipmentCreationDTO represents the data required to create a new equipment entry
type EquipmentCreationDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
