package dto

import (
	"time"

	"github.com/lib/pq"
)

type GymResponseDTO struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Domain         string         `json:"domain"`
	Email          string         `json:"email"`
	Address        string         `json:"address"`
	Phone          string         `json:"phone"`
	IsActive       bool           `json:"is_active"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	BusinessHours  pq.StringArray `json:"business_hours"`
	SocialLinks    pq.StringArray `json:"social_links"`
	PaymentMethods pq.StringArray `json:"payment_methods"`
}
