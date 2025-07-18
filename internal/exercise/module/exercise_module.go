package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise/handler"
	"github.com/alejandro-albiol/athenai/internal/exercise/repository"
	"github.com/alejandro-albiol/athenai/internal/exercise/router"
	"github.com/alejandro-albiol/athenai/internal/exercise/service"
)

func NewExerciseModule(db *sql.DB) http.Handler {
	repo := repository.NewExerciseRepository(db)
	service := service.NewExerciseService(repo)
	handler := handler.NewExerciseHandler(service)
	return router.NewExerciseRouter(handler)
}
