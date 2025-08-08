package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_template/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewWorkoutTemplateRouter(handler interfaces.WorkoutTemplateHandler) http.Handler {
	r := chi.NewRouter()

	// Workout template CRUD endpoints
	r.Post("/", handler.CreateWorkoutTemplate)                                              // POST /workout-templates
	r.Get("/", handler.GetAllWorkoutTemplates)                                              // GET /workout-templates
	r.Get("/{id}", handler.GetWorkoutTemplateByID)                                          // GET /workout-templates/{id}
	r.Get("/name/{name}", handler.GetWorkoutTemplateByName)                                 // GET /workout-templates/name/{name}
	r.Get("/difficulty/{difficulty}", handler.GetWorkoutTemplatesByDifficulty)              // GET /workout-templates/difficulty/{difficulty}
	r.Get("/target-audience/{targetAudience}", handler.GetWorkoutTemplatesByTargetAudience) // GET /workout-templates/target-audience/{targetAudience}
	r.Put("/{id}", handler.UpdateWorkoutTemplate)                                           // PUT /workout-templates/{id}
	r.Delete("/{id}", handler.DeleteWorkoutTemplate)                                        // DELETE /workout-templates/{id}

	return r
}
