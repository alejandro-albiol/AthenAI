package interfaces

import (
	"github.com/alejandro-albiol/athenai/internal/workout_generator/dto"
)

type WorkoutGeneratorService interface {
	GenerateWorkout(req *dto.WorkoutGeneratorRequest) (*dto.WorkoutGeneratorResponse, error)
}
