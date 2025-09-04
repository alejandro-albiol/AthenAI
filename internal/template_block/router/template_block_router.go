package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewTemplateBlockRouter(h interfaces.TemplateBlockHandlerInterface) http.Handler {
	r := chi.NewRouter()
	r.Post("/template-blocks", h.CreateTemplateBlock)
	r.Get("/template-blocks/{id}", h.GetTemplateBlockByID)
	r.Put("/template-blocks/{id}", h.UpdateTemplateBlock)
	r.Delete("/template-blocks/{id}", h.DeleteTemplateBlock)
	r.Get("/templates/{templateId}/blocks", h.ListTemplateBlocksByTemplateID)
	return r
}
