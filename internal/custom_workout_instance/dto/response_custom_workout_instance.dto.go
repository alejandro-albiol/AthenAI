package dto

import "github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"

type ResponseCustomWorkoutInstanceDTO struct {
	// Basic Information
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	TemplateSource   string  `json:"template_source"`
	PublicTemplateID *string `json:"public_template_id,omitempty"`
	GymTemplateID    *string `json:"gym_template_id,omitempty"`

	// Calculated Fields (computed from exercises)
	DifficultyLevel          string `json:"difficulty_level"`           // Calculated from average exercise difficulty
	EstimatedDurationMinutes int    `json:"estimated_duration_minutes"` // Calculated from exercises + rest times
	TotalExercises           int    `json:"total_exercises"`            // Count of exercises in this workout
	TotalSets                int    `json:"total_sets"`                 // Sum of all sets across exercises

	// Exercise Type Breakdown
	ExerciseTypes   []string `json:"exercise_types"`   // Unique exercise types (strength, cardio, etc.)
	MuscularGroups  []string `json:"muscular_groups"`  // Unique muscular groups targeted
	EquipmentNeeded []string `json:"equipment_needed"` // All equipment required

	// Workout Statistics
	WorkoutStats *WorkoutStatsDTO `json:"workout_stats,omitempty"`

	// Optional: Include exercises for detailed view
	Exercises []dto.ResponseCustomWorkoutExerciseDTO `json:"exercises,omitempty"`

	// Timestamps
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type WorkoutStatsDTO struct {
	AverageSetsPerExercise float64        `json:"average_sets_per_exercise"`
	AverageRepsPerSet      float64        `json:"average_reps_per_set"`
	TotalEstimatedWeight   float64        `json:"total_estimated_weight"`  // Sum of all weights
	DifficultyBreakdown    map[string]int `json:"difficulty_breakdown"`    // {"beginner": 2, "intermediate": 3}
	ExerciseTypeBreakdown  map[string]int `json:"exercise_type_breakdown"` // {"strength": 5, "cardio": 2}
}
