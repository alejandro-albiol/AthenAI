package dto

type UpdateCustomMemberWorkoutDTO struct {
	ID          string  `json:"id"`
	StartedAt   *string `json:"started_at,omitempty"`
	CompletedAt *string `json:"completed_at,omitempty"`
	Status      *string `json:"status,omitempty"`
	Notes       *string `json:"notes,omitempty"`
	Rating      *int    `json:"rating,omitempty"`
}
