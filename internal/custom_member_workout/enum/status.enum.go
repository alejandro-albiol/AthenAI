package enum

type CustomMemberWorkoutStatus string

const (
	StatusScheduled  CustomMemberWorkoutStatus = "scheduled"
	StatusInProgress CustomMemberWorkoutStatus = "in_progress"
	StatusCompleted  CustomMemberWorkoutStatus = "completed"
	StatusSkipped    CustomMemberWorkoutStatus = "skipped"
	StatusCancelled  CustomMemberWorkoutStatus = "cancelled"
)

func (e CustomMemberWorkoutStatus) IsValid() bool {
	switch e {
	case StatusScheduled, StatusInProgress, StatusCompleted, StatusSkipped, StatusCancelled:
		return true
	}
	return false
}
