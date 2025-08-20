package enum

type DifficultyLevel string

const (
	Beginner     DifficultyLevel = "beginner"
	Intermediate DifficultyLevel = "intermediate"
	Advanced     DifficultyLevel = "advanced"
)

func (d DifficultyLevel) IsValid() bool {
	switch d {
	case Beginner, Intermediate, Advanced:
		return true
	}
	return false
}
