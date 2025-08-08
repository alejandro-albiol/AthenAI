package interfaces

import "github.com/alejandro-albiol/athenai/internal/workout_template/dto"

// WorkoutTemplateService defines service methods for workout templates.
type WorkoutTemplateService interface {
	// CreateWorkoutTemplate creates a new workout template in the system.
	CreateWorkoutTemplate(dto.CreateWorkoutTemplateDTO) (string, error)

	// GetWorkoutTemplateByID retrieves a workout template by its unique ID.
	GetWorkoutTemplateByID(id string) (dto.WorkoutTemplateDTO, error)

	// GetWorkoutTemplateByName retrieves a workout template by its name.
	GetWorkoutTemplateByName(name string) (dto.WorkoutTemplateDTO, error)

	// GetWorkoutTemplatesByDifficulty retrieves all workout templates matching a given difficulty level.
	GetWorkoutTemplatesByDifficulty(difficulty string) ([]dto.WorkoutTemplateDTO, error)

	// GetWorkoutTemplatesByTargetAudience retrieves all workout templates for a specific target audience.
	GetWorkoutTemplatesByTargetAudience(targetAudience string) ([]dto.WorkoutTemplateDTO, error)

	// GetAllWorkoutTemplates retrieves all workout templates in the system.
	GetAllWorkoutTemplates() ([]dto.WorkoutTemplateDTO, error)

	// UpdateWorkoutTemplate updates an existing workout template by ID.
	UpdateWorkoutTemplate(id string, dto dto.UpdateWorkoutTemplateDTO) (dto.WorkoutTemplateDTO, error)

	// DeleteWorkoutTemplate deletes a workout template by ID.
	DeleteWorkoutTemplate(id string) error
}
