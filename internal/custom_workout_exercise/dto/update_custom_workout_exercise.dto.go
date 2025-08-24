package dto

type UpdateCustomWorkoutExerciseDTO struct {
	ID              string   `json:"id"`
	Sets            *int     `json:"sets,omitempty"`
	RepsMin         *int     `json:"reps_min,omitempty"`
	RepsMax         *int     `json:"reps_max,omitempty"`
	WeightKg        *float64 `json:"weight_kg,omitempty"`
	DurationSeconds *int     `json:"duration_seconds,omitempty"`
	RestSeconds     *int     `json:"rest_seconds,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}
