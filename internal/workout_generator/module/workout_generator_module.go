package module

import (
	"net/http"
	"os"

	exerciseIF "github.com/alejandro-albiol/athenai/internal/exercise/interfaces"
	templateIF "github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
	userIF "github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/handler"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/router"
	"github.com/alejandro-albiol/athenai/internal/workout_generator/service"
	workoutTemplateIF "github.com/alejandro-albiol/athenai/internal/workout_template/interfaces"
)

// NewWorkoutGeneratorModule wires up repository, service, handler, router for a gym
func NewWorkoutGeneratorModule(
	exerciseSvc exerciseIF.ExerciseService,
	workoutTemplateSvc workoutTemplateIF.WorkoutTemplateService,
	templateBlockSvc templateIF.TemplateBlockService,
	userSvc userIF.UserService,
) http.Handler {
	svc := service.NewWorkoutGeneratorService(
		os.Getenv("LLM_ENDPOINT"),
		os.Getenv("API_TOKEN"),
		exerciseSvc,
		workoutTemplateSvc,
		templateBlockSvc,
		userSvc,
	)
	h := handler.NewWorkoutGeneratorHandler(svc)
	return router.NewWorkoutGeneratorRouter(h)
}
