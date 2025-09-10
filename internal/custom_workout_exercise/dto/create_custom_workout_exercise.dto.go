package dto

type CreateCustomWorkoutExerciseDTO struct {
	CreatedBy         string  `json:"created_by" validate:"required"`
	WorkoutInstanceID string  `json:"workout_instance_id" validate:"required"`
	ExerciseSource    string  `json:"exercise_source" validate:"required,oneof=public gym"`
	PublicExerciseID  *string `json:"public_exercise_id,omitempty"`
	GymExerciseID     *string `json:"gym_exercise_id,omitempty"`
	BlockName         string  `json:"block_name" validate:"required"`
	ExerciseOrder     int     `json:"exercise_order" validate:"required,min=1"`
	// Actual execution parameters (can override template block defaults)
	Sets            *int     `json:"sets,omitempty" validate:"omitempty,min=1"`
	RepsMin         *int     `json:"reps_min,omitempty" validate:"omitempty,min=1"`
	RepsMax         *int     `json:"reps_max,omitempty" validate:"omitempty,min=1"`
	WeightKg        *float64 `json:"weight_kg,omitempty" validate:"omitempty,min=0"`
	DurationSeconds *int     `json:"duration_seconds,omitempty" validate:"omitempty,min=1"`
	RestSeconds     *int     `json:"rest_seconds,omitempty" validate:"omitempty,min=0"`
	Notes           *string  `json:"notes,omitempty"`
}
