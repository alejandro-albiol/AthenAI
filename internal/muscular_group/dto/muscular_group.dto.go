package dto

// CreateMuscularGroupDTO represents the data required to create a new muscular group
type CreateMuscularGroupDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BodyPart    string `json:"body_part"` // 'upper_body', 'lower_body', 'core', 'full_body'
}

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

// UpdateMuscularGroupDTO represents the data for updating a muscular group
type UpdateMuscularGroupDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	BodyPart    *string `json:"body_part,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
