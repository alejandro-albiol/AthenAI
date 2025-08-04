package enums

// DifficultyLevel represents the different difficulty levels in the application.
type DifficultyLevel string

// List of all difficulty levels.
const (
	Beginner     DifficultyLevel = "beginner"
	Intermediate DifficultyLevel = "intermediate"
	Advanced     DifficultyLevel = "advanced"
)

// AllDifficultyLevels is a slice of all valid difficulty level values for validation.
var AllDifficultyLevels = []DifficultyLevel{
	Beginner, Intermediate, Advanced,
}
