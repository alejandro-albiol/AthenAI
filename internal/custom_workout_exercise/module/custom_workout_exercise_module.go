package module

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/router"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/service"

	"net/http"
)

func NewCustomWorkoutExerciseModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomWorkoutExerciseRepository(db)
	service := service.NewCustomWorkoutExerciseService(repo)
	handler := handler.NewCustomWorkoutExerciseHandler(service)
	return router.NewCustomWorkoutExerciseRouter(handler)
}
