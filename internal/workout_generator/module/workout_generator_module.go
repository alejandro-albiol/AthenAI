package module

import (
	"net/http"
	"os"

	"github.com/alejandro-albiol/athenai/internal/workout_generator/handler"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/router"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/service"
)

// NewWorkoutGeneratorModule wires up repository, service, handler, router for a gym
func NewWorkoutGeneratorModule() http.Handler {

	svc := service.NewWorkoutGeneratorService(os.Getenv("LLM_ENDPOINT"), os.Getenv("API_TOKEN"))
	h := handler.NewWorkoutGeneratorHandler(svc)
	return router.NewWorkoutGeneratorRouter(h)
}
