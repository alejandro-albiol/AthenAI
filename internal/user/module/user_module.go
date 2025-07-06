package usermodule

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/handler"
	"github.com/alejandro-albiol/athenai/internal/user/repository"
	"github.com/alejandro-albiol/athenai/internal/user/router"
	"github.com/alejandro-albiol/athenai/internal/user/service"
)

func NewUserModule(db *sql.DB) http.Handler {
	repo := repository.NewUsersRepository(db)
	service := service.NewUsersService(repo)
	handler := handler.NewUsersHandler(service)
	return router.NewUsersRouter(handler)
}
