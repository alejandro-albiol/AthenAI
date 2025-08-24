package dto

type CreateCustomMemberWorkoutDTO struct {
	MemberID          string  `json:"member_id"`
	WorkoutInstanceID string  `json:"workout_instance_id"`
	ScheduledDate     string  `json:"scheduled_date"`
	Notes             *string `json:"notes,omitempty"`
	Rating            *int    `json:"rating,omitempty"`
}
