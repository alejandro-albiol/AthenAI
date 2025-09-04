package dto

// UpdateMuscularGroupDTO represents the data for updating a muscular group
type UpdateMuscularGroupDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	BodyPart    *string `json:"body_part,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
