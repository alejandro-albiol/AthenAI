package usermodule

import (
	"database/sql"
	"net/http"

	gymrepository "github.com/alejandro-albiol/athenai/internal/gym/repository"
	"github.com/alejandro-albiol/athenai/internal/user/handler"
	"github.com/alejandro-albiol/athenai/internal/user/repository"
	"github.com/alejandro-albiol/athenai/internal/user/router"
	"github.com/alejandro-albiol/athenai/internal/user/service"
)

func NewUserModule(db *sql.DB) http.Handler {
	// Create gym repository for dependency injection
	gymRepo := gymrepository.NewGymRepository(db)

	// Create user repository with gym repository dependency
	repo := repository.NewUsersRepository(db, gymRepo)
	service := service.NewUsersService(repo)
	handler := handler.NewUsersHandler(service)
	return router.NewUsersRouter(handler)
}
