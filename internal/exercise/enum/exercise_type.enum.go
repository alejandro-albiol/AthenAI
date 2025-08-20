package enum

type ExerciseType string

const (
	Strength    ExerciseType = "strength"
	Cardio      ExerciseType = "cardio"
	Flexibility ExerciseType = "flexibility"
	Balance     ExerciseType = "balance"
	Functional  ExerciseType = "functional"
)

func (e ExerciseType) IsValid() bool {
	switch e {
	case Strength, Cardio, Flexibility, Balance, Functional:
		return true
	}
	return false
}
