package dto

type GymUpdateDTO struct {
	Name        string   `json:"name,omitempty" validate:"omitempty"`
	Domain      string   `json:"domain,omitempty" validate:"omitempty"`
	Email       string   `json:"email,omitempty" validate:"omitempty,email"`
	Address     string   `json:"address,omitempty" validate:"omitempty"`
	ContactName string   `json:"contact_name,omitempty" validate:"omitempty"`
	Phone       string   `json:"phone,omitempty" validate:"omitempty"`
	LogoURL     string   `json:"logo_url,omitempty" validate:"omitempty"`
}