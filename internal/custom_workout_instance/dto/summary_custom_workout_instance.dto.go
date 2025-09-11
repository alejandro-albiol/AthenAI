package dto

type SummaryCustomWorkoutInstanceDTO struct {
	// Basic Information
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	TemplateSource string `json:"template_source"`

	// Key Metrics
	DifficultyLevel          string `json:"difficulty_level"`
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes"`
	TotalExercises           int    `json:"total_exercises"`
	TotalSets                int    `json:"total_sets"`

	// Primary Exercise Types and Muscle Groups
	PrimaryExerciseType   string   `json:"primary_exercise_type"`   // Most common exercise type
	PrimaryMuscularGroups []string `json:"primary_muscular_groups"` // Top 3 muscle groups

	// Timestamps
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
