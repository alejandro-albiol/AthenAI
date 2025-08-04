package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/template_block/handler"
	"github.com/alejandro-albiol/athenai/internal/template_block/repository"
	"github.com/alejandro-albiol/athenai/internal/template_block/router"
	"github.com/alejandro-albiol/athenai/internal/template_block/service"
)

func NewTemplateBlockModule(db *sql.DB) http.Handler {
	repo := repository.NewTemplateBlockRepository(db)
	service := service.NewTemplateBlockService(repo)
	handler := handler.NewTemplateBlockHandler(service)
	return router.NewTemplateBlockRouter(handler)
}
