package dto

type ResponseCustomMemberWorkoutDTO struct {
	ID                string  `json:"id"`
	MemberID          string  `json:"member_id"`
	WorkoutInstanceID string  `json:"workout_instance_id"`
	ScheduledDate     string  `json:"scheduled_date"`
	StartedAt         *string `json:"started_at,omitempty"`
	CompletedAt       *string `json:"completed_at,omitempty"`
	Status            string  `json:"status"`
	Notes             *string `json:"notes,omitempty"`
	Rating            *int    `json:"rating,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}
