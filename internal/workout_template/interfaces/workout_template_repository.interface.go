package interfaces

import "github.com/alejandro-albiol/athenai/internal/workout_template/dto"

type WorkoutTemplateRepository interface {
	// Create inserts a new workout template into the database.
	CreateWorkoutTemplate(dto *dto.CreateWorkoutTemplateDTO) (string, error)

	// GetByID fetches a workout template by its unique ID.
	GetWorkoutTemplateByID(id string) (*dto.ResponseWorkoutTemplateDTO, error)

	// GetByName fetches a workout template by its name.
	GetWorkoutTemplateByName(name string) (*dto.ResponseWorkoutTemplateDTO, error)

	// GetByDifficulty fetches all workout templates with the specified difficulty.
	GetWorkoutTemplatesByDifficulty(difficulty string) ([]*dto.ResponseWorkoutTemplateDTO, error)

	// GetByTargetAudience fetches all workout templates for a given target audience.
	GetWorkoutTemplatesByTargetAudience(targetAudience string) ([]*dto.ResponseWorkoutTemplateDTO, error)

	// GetAll fetches all workout templates from the database.
	GetAllWorkoutTemplates() ([]*dto.ResponseWorkoutTemplateDTO, error)

	// Update modifies an existing workout template by ID.
	UpdateWorkoutTemplate(id string, dto *dto.UpdateWorkoutTemplateDTO) (*dto.ResponseWorkoutTemplateDTO, error)

	// Delete removes a workout template by ID.
	DeleteWorkoutTemplate(id string) error
}
