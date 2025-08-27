package dto

type ResponseCustomEquipmentDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CreatedBy   string `json:"created_by"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
