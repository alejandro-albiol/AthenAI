package dto

type UpdateCustomWorkoutExerciseDTO struct {
	ID              string   `json:"id" validate:"required"`
	Sets            *int     `json:"sets,omitempty" validate:"omitempty,min=1"`
	RepsMin         *int     `json:"reps_min,omitempty" validate:"omitempty,min=1"`
	RepsMax         *int     `json:"reps_max,omitempty" validate:"omitempty,min=1"`
	WeightKg        *float64 `json:"weight_kg,omitempty" validate:"omitempty,min=0"`
	DurationSeconds *int     `json:"duration_seconds,omitempty" validate:"omitempty,min=1"`
	RestSeconds     *int     `json:"rest_seconds,omitempty" validate:"omitempty,min=0"`
	Notes           *string  `json:"notes,omitempty"`
}
