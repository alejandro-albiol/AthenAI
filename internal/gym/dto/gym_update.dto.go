package dto

type GymUpdateDTO struct {
	Name    *string `json:"name,omitempty" validate:"omitempty"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
	Address *string `json:"address,omitempty" validate:"omitempty"`
	Phone   *string `json:"phone,omitempty" validate:"omitempty"`
}
