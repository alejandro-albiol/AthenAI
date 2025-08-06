package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/workout_template/handler"
	"github.com/alejandro-albiol/athenai/internal/workout_template/repository"
	"github.com/alejandro-albiol/athenai/internal/workout_template/router"
	"github.com/alejandro-albiol/athenai/internal/workout_template/service"
)

func NewWorkoutTemplateModule(db *sql.DB) http.Handler {
	repo := repository.NewWorkoutTemplateRepository(db)
	service := service.NewWorkoutTemplateService(repo)
	handler := handler.NewWorkoutTemplateHandler(service)
	return router.NewWorkoutTemplateRouter(handler)
}
