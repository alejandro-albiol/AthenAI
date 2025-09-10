package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/interfaces"
)

func NewCustomTemplateBlockRouter(h interfaces.CustomTemplateBlockHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/custom-template-block", h.Create)
	r.Get("/custom-template-block/{id}", h.GetByID)
	r.Get("/custom-template-block/template/{templateID}", h.ListByTemplateID)
	r.Put("/custom-template-block/{id}", h.Update)
	r.Delete("/custom-template-block/{id}", h.Delete)
	return r
}
