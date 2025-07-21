package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/handler"
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/repository"
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/router"
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/service"
)

func NewExerciseMuscularGroupModule(db *sql.DB) http.Handler {
	repo := repository.NewExerciseMuscularGroupRepository(db)
	service := service.NewExerciseMuscularGroupService(repo)
	handler := handler.NewExerciseMuscularGroupHandler(service)
	return router.NewExerciseMuscularGroupRouter(handler)
}
