package dto

// InvitationDecodeResponseDTO - Response for decoding invitation tokens
type InvitationDecodeResponseDTO struct {
	GymID   string `json:"gym_id"`
	GymName string `json:"gym_name"`
	Valid   bool   `json:"valid"`
}
