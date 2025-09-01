package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/router"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/service"
	muscular_group_repository "github.com/alejandro-albiol/athenai/internal/muscular_group/repository"
)

func NewCustomExerciseMuscularGroupModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomExerciseMuscularGroupRepository(db)
	publicMuscularGroupRepo := muscular_group_repository.NewMuscularGroupRepository(db)
	svc := service.NewCustomExerciseMuscularGroupService(repo, publicMuscularGroupRepo)
	handler := handler.NewCustomExerciseMuscularGroupHandler(svc)
	router := router.NewCustomExerciseMuscularGroupRouter(handler)
	return router
}
