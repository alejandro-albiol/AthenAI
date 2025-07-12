package dto

import "time"

type MuscularGroupCreationDTO struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	BodyPart    string  `json:"body_part" validate:"required,oneof=upper_body lower_body core full_body"`
}

type MuscularGroupUpdateDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	BodyPart    *string `json:"body_part" validate:"omitempty,oneof=upper_body lower_body core full_body"`
	IsActive    *bool   `json:"is_active"`
}

type MuscularGroupResponseDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	BodyPart    string    `json:"body_part"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
