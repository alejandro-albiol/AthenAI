package dto

// MuscularGroupResponseDTO represents a muscular group in API responses
type MuscularGroupResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BodyPart    string `json:"body_part"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
