package dto

import "time"

// InvitationCreateRequestDTO - Request for creating gym invitations
type InvitationCreateRequestDTO struct {
	Email        string `json:"email" validate:"required,email"`
	Role         string `json:"role" validate:"required,oneof=gym_admin trainer member"`
	GymID        string `json:"gym_id" validate:"required"`
	ExpiresHours int    `json:"expires_hours,omitempty"` // Default 72 hours
	Message      string `json:"message,omitempty"`       // Optional personal message
}

// InvitationResponseDTO - Response for invitation operations
type InvitationResponseDTO struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	GymID      string     `json:"gym_id"`
	GymName    string     `json:"gym_name"`
	Token      string     `json:"token"`
	InviteURL  string     `json:"invite_url"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	CreatedBy  string     `json:"created_by"`
}

// InvitationListResponseDTO - Response for listing invitations
type InvitationListResponseDTO struct {
	Invitations []InvitationResponseDTO `json:"invitations"`
	Total       int                     `json:"total"`
}

// InvitationDecodeResponseDTO - Response for decoding invitation tokens
type InvitationDecodeResponseDTO struct {
	GymID       string `json:"gym_id"`
	GymName     string `json:"gym_name"`
	Description string `json:"description,omitempty"`
	Address     string `json:"address,omitempty"`
	Role        string `json:"role"`
	Email       string `json:"email"`
	Valid       bool   `json:"valid"`
	Expired     bool   `json:"expired"`
}

// InvitationAcceptRequestDTO - Request for accepting invitations
type InvitationAcceptRequestDTO struct {
	Token     string `json:"token" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
