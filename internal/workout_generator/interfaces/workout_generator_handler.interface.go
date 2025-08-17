package interfaces

// Handler interface for workout generator HTTP endpoints
// Follows AthenAI conventions: handler only translates service results to HTTP responses
// Service layer contains all business logic

import (
	"net/http"
)

type WorkoutGeneratorHandler interface {
	// GenerateWorkout handles POST /workout/generate
	GenerateWorkout(w http.ResponseWriter, r *http.Request)
}
