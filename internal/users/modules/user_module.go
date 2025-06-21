package usersmodule

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/databases/interfaces"
	"github.com/alejandro-albiol/athenai/internal/users/handlers"
	"github.com/alejandro-albiol/athenai/internal/users/repositories"
	"github.com/alejandro-albiol/athenai/internal/users/routers"
	"github.com/alejandro-albiol/athenai/internal/users/services"
)

func NewUserModule(db interfaces.DBService) http.Handler {
    repo := repository.NewUsersRepository(db)
    service := services.NewUsersService(repo)
    handler := handlers.NewUsersHandler(service)
    return router.NewUsersRouter(handler)
}