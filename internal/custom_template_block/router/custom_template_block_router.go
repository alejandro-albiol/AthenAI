package router

import (
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/handler"
	"github.com/go-chi/chi/v5"
)

func NewCustomTemplateBlockRouter(h *handler.CustomTemplateBlockHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/custom-template-block", h.Create)
	r.Get("/custom-template-block/{id}", h.GetByID)
	r.Get("/custom-template-block/template/{templateID}", h.ListByTemplateID)
	r.Put("/custom-template-block/{id}", h.Update)
	r.Delete("/custom-template-block/{id}", h.Delete)
	return r
}
