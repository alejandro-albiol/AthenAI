package interfaces

import (
	"net/http"
)

type WorkoutTemplateHandler interface {
	// CreateWorkoutTemplate handles HTTP requests to create a new workout template.
	CreateWorkoutTemplate(w http.ResponseWriter, r *http.Request)

	// GetWorkoutTemplateByID handles HTTP requests to retrieve a workout template by ID.
	GetWorkoutTemplateByID(w http.ResponseWriter, r *http.Request)

	// GetWorkoutTemplateByName handles HTTP requests to retrieve a workout template by name.
	GetWorkoutTemplateByName(w http.ResponseWriter, r *http.Request)

	// GetWorkoutTemplatesByDifficulty handles HTTP requests to retrieve workout templates by difficulty.
	GetWorkoutTemplatesByDifficulty(w http.ResponseWriter, r *http.Request)

	// GetWorkoutTemplatesByTargetAudience handles HTTP requests to retrieve workout templates by target audience.
	GetWorkoutTemplatesByTargetAudience(w http.ResponseWriter, r *http.Request)

	// GetAllWorkoutTemplates handles HTTP requests to retrieve all workout templates.
	GetAllWorkoutTemplates(w http.ResponseWriter, r *http.Request)

	// UpdateWorkoutTemplate handles HTTP requests to update a workout template.
	UpdateWorkoutTemplate(w http.ResponseWriter, r *http.Request)

	// DeleteWorkoutTemplate handles HTTP requests to delete a workout template.
	DeleteWorkoutTemplate(w http.ResponseWriter, r *http.Request)
}
