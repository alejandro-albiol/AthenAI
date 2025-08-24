package dto

type CreateCustomWorkoutExerciseDTO struct {
	WorkoutInstanceID string   `json:"workout_instance_id"`
	ExerciseSource    string   `json:"exercise_source"`
	PublicExerciseID  *string  `json:"public_exercise_id,omitempty"`
	GymExerciseID     *string  `json:"gym_exercise_id,omitempty"`
	BlockName         string   `json:"block_name"`
	ExerciseOrder     int      `json:"exercise_order"`
	Sets              *int     `json:"sets,omitempty"`
	RepsMin           *int     `json:"reps_min,omitempty"`
	RepsMax           *int     `json:"reps_max,omitempty"`
	WeightKg          *float64 `json:"weight_kg,omitempty"`
	DurationSeconds   *int     `json:"duration_seconds,omitempty"`
	RestSeconds       *int     `json:"rest_seconds,omitempty"`
	Notes             *string  `json:"notes,omitempty"`
}
