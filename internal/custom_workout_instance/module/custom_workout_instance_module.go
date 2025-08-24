package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/service"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/router"
)

func NewCustomWorkoutInstanceModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomWorkoutInstanceRepository(db)
	service := service.NewCustomWorkoutInstanceService(repo)
	handler := handler.NewCustomWorkoutInstanceHandler(service)
	return router.NewCustomWorkoutInstanceRouter(handler)
}
