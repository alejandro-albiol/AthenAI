package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/muscular_group/handler"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/repository"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/router"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/service"
)

func NewMuscularGroupModule(db *sql.DB) http.Handler {
	repo := repository.NewMuscularGroupRepository(db)
	service := service.NewMuscularGroupService(repo)
	handler := handler.NewMuscularGroupHandler(service)
	return router.NewMuscularGroupRouter(handler)
}
