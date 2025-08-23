package router

import (
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/handler"
	"github.com/go-chi/chi/v5"
)

func NewCustomEquipmentRouter(h *handler.CustomEquipmentHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/custom-equipment", h.Create)
	r.Get("/custom-equipment/{id}", h.GetByID)
	r.Get("/custom-equipment", h.List)
	r.Put("/custom-equipment/{id}", h.Update)
	r.Delete("/custom-equipment/{id}", h.Delete)
	return r
}
