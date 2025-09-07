package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewTemplateBlockRouter(h interfaces.TemplateBlockHandlerInterface) http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.CreateTemplateBlock)
	r.Get("/{id}", h.GetTemplateBlockByID)
	r.Put("/{id}", h.UpdateTemplateBlock)
	r.Delete("/{id}", h.DeleteTemplateBlock)
	r.Get("/template/{templateId}/block", h.ListTemplateBlocksByTemplateID)
	return r
}
