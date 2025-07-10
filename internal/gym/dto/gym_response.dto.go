package dto

import "time"

type GymResponseDTO struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Domain         string    `json:"domain"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`
	ContactName    string    `json:"contact_name"`
	Phone          string    `json:"phone"`
	LogoURL        string    `json:"logo_url,omitempty"`
	IsActive       bool      `json:"is_active"`
	Description    string    `json:"description,omitempty"`
	BusinessHours  []string  `json:"business_hours,omitempty"`
	SocialLinks    []string  `json:"social_links,omitempty"`
	PaymentMethods []string  `json:"payment_methods,omitempty"`
	Currency       string    `json:"currency,omitempty"`
	TimezoneOffset string    `json:"timezone_offset,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
