package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/gym/handler"
	"github.com/alejandro-albiol/athenai/internal/gym/repository"
	"github.com/alejandro-albiol/athenai/internal/gym/router"
	"github.com/alejandro-albiol/athenai/internal/gym/service"
)

func NewGymModule(db *sql.DB) http.Handler {
	repo := repository.NewGymRepository(db)
	service := service.NewGymService(repo)
	handler := handler.NewGymHandler(service)
	return router.NewGymRouter(handler)
}
