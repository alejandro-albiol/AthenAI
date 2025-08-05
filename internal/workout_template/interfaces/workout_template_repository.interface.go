package interfaces

import "github.com/alejandro-albiol/athenai/internal/workout_template/dto"

type WorkoutTemplateRepository interface {
	// Create inserts a new workout template into the database.
	Create(template dto.CreateWorkoutTemplateDTO) error

	// GetByID fetches a workout template by its unique ID.
	GetByID(id string) (dto.WorkoutTemplateDTO, error)

	// GetByName fetches a workout template by its name.
	GetByName(name string) (dto.WorkoutTemplateDTO, error)

	// GetByDifficulty fetches all workout templates with the specified difficulty.
	GetByDifficulty(difficulty string) ([]dto.WorkoutTemplateDTO, error)

	// GetByTargetAudience fetches all workout templates for a given target audience.
	GetByTargetAudience(targetAudience string) ([]dto.WorkoutTemplateDTO, error)

	// GetAll fetches all workout templates from the database.
	GetAll() ([]dto.WorkoutTemplateDTO, error)

	// Update modifies an existing workout template by ID.
	Update(id string, template dto.UpdateWorkoutTemplateDTO) (dto.WorkoutTemplateDTO, error)

	// Delete removes a workout template by ID.
	Delete(id string) error
}
