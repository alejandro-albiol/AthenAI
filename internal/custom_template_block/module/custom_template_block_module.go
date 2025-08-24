package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/router"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/service"
)

func NewCustomTemplateBlockModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomTemplateBlockRepository(db)
	service := service.NewCustomTemplateBlockService(repo)
	handler := handler.NewCustomTemplateBlockHandler(service)
	return router.NewCustomTemplateBlockRouter(handler)
}
