package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewTemplateBlockRouter(h interfaces.TemplateBlockHandlerInterface) http.Handler {
	r := chi.NewRouter()
	r.Post("/template-block", h.CreateTemplateBlock)
	r.Get("/template-block/{id}", h.GetTemplateBlockByID)
	r.Put("/template-block/{id}", h.UpdateTemplateBlock)
	r.Delete("/template-block/{id}", h.DeleteTemplateBlock)
	r.Get("/template/{templateId}/block", h.ListTemplateBlocksByTemplateID)
	return r
}
