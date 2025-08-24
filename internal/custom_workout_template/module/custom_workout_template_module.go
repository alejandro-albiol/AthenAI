package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/router"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/service"
)

func NewCustomWorkoutTemplateModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomWorkoutTemplateRepository(db)
	service := service.NewCustomWorkoutTemplateService(repo)
	handler := handler.NewCustomWorkoutTemplateHandler(service)
	return router.NewCustomWorkoutTemplateRouter(handler)
}
