package dto

type CreateCustomEquipmentDTO struct {
	Name        string `json:"name"`
	CreatedBy   string `json:"created_by"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsActive    bool   `json:"is_active"`
}
