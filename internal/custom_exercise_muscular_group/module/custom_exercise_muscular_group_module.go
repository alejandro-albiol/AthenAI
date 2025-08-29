package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/router"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/service"
)

func NewCustomExerciseMuscularGroupModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomExerciseMuscularGroupRepository(db)
	svc := service.NewCustomExerciseMuscularGroupService(repo)
	handler := handler.NewCustomExerciseMuscularGroupHandler(svc)
	router := router.NewCustomExerciseMuscularGroupRouter(handler)
	return router
}
