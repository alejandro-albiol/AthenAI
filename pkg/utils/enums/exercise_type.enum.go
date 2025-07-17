package enums

// ExerciseType represents the different types of exercises in the application.
type ExerciseType string

// List of all exercise types.
const (
	Strength    ExerciseType = "strength"
	Cardio      ExerciseType = "cardio"
	Flexibility ExerciseType = "flexibility"
	Balance     ExerciseType = "balance"
	Functional  ExerciseType = "functional"
)

// AllExerciseTypes is a slice of all valid exercise type values for validation.
var AllExerciseTypes = []ExerciseType{
	Strength, Cardio, Flexibility, Balance, Functional,
}
