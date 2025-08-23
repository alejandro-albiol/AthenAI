package dto

type UpdateCustomEquipmentDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
