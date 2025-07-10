package dto

type GymCreationDTO struct {
	Name        string   `json:"name" validate:"required"`
	Domain      string   `json:"domain" validate:"required"`
	Email       string   `json:"email" validate:"required,email"`
	Address     string   `json:"address" validate:"required"`
	ContactName string   `json:"contact_name" validate:"required"`
	Phone       string   `json:"phone" validate:"required"`
	LogoURL     string   `json:"logo_url,omitempty"`
}