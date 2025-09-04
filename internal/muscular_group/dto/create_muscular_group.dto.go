package dto

// CreateMuscularGroupDTO represents the data required to create a new muscular group
type CreateMuscularGroupDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BodyPart    string `json:"body_part"` // 'upper_body', 'lower_body', 'core', 'full_body'
}
