package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/router"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/service"
)

// NewCustomEquipmentModule wires up the custom equipment module
func NewCustomEquipmentModule(db *sql.DB) http.Handler {
	repo := &repository.CustomEquipmentRepositoryImpl{DB: db}
	service := &service.CustomEquipmentServiceImpl{Repo: repo}
	handler := &handler.CustomEquipmentHandler{Service: service}
	return router.NewCustomEquipmentRouter(handler)
}
