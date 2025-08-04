package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/template_block/handler"
	"github.com/go-chi/chi/v5"
)

func NewTemplateBlockRouter(h *handler.TemplateBlockHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/template-blocks", h.CreateTemplateBlock)
	r.Get("/template-blocks/{id}", h.GetTemplateBlock)
	r.Put("/template-blocks/{id}", h.UpdateTemplateBlock)
	r.Delete("/template-blocks/{id}", h.DeleteTemplateBlock)
	return r
}
