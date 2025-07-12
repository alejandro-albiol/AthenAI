package dto

import "time"

type EquipmentCreationDTO struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Category    string  `json:"category" validate:"required,oneof=free_weights machines cardio accessories bodyweight"`
}

type EquipmentUpdateDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category" validate:"omitempty,oneof=free_weights machines cardio accessories bodyweight"`
	IsActive    *bool   `json:"is_active"`
}

type EquipmentResponseDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Category    string    `json:"category"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
