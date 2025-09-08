package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewCustomEquipmentRouter(h interfaces.CustomEquipmentHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/custom-equipment", h.CreateCustomEquipment)
	r.Get("/custom-equipment/{id}", h.GetCustomEquipmentByID)
	r.Get("/custom-equipment", h.ListCustomEquipment)
	r.Get("/custom-equipment/search", h.GetCustomEquipmentByName)
	r.Put("/custom-equipment/{id}", h.UpdateCustomEquipment)
	r.Delete("/custom-equipment/{id}", h.DeleteCustomEquipment)
	return r
}
