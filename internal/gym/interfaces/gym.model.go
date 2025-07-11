package interfaces

import "time"

type Gym struct {
	ID             string     `json:"id"`
	Name           string     `json:"name" validate:"required"`
	Domain         string     `json:"domain" validate:"required"`
	Email          string     `json:"email" validate:"required,email"`
	Address        string     `json:"address" validate:"required"`
	Phone          string     `json:"phone" validate:"required"`
	IsActive       bool       `json:"is_active"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	BusinessHours  []string   `json:"business_hours,omitempty"`
	SocialLinks    []string   `json:"social_links,omitempty"`
	PaymentMethods []string   `json:"payment_methods,omitempty"`
}
