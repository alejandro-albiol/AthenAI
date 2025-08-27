package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/router"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/service"
)

type CustomExerciseModule struct {
	Handler *handler.CustomExerciseHandler
	Router  http.Handler
}

func NewCustomExerciseModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomExerciseRepository(db)
	svc := service.NewCustomExerciseService(repo)
	handler := handler.NewCustomExerciseHandler(svc)
	return router.NewCustomExerciseRouter(handler)
}
