package dto

import (
	"time"
)

type GymResponseDTO struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Domain         string     `json:"domain"`
	Email          string     `json:"email"`
	Address        string     `json:"address"`
	Phone          string     `json:"phone"`
	IsActive       bool       `json:"is_active"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	BusinessHours  []string   `json:"business_hours"`
	SocialLinks    []string   `json:"social_links"`
	PaymentMethods []string   `json:"payment_methods"`
}
