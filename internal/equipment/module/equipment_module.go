package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/equipment/handler"
	"github.com/alejandro-albiol/athenai/internal/equipment/repository"
	"github.com/alejandro-albiol/athenai/internal/equipment/router"
	"github.com/alejandro-albiol/athenai/internal/equipment/service"
)

func NewEquipmentModule(db *sql.DB) http.Handler {
	repo := repository.NewEquipmentRepository(db)
	service := service.NewEquipmentService(repo)
	handler := handler.NewEquipmentHandler(service)
	return router.NewEquipmentRouter(handler)
}
