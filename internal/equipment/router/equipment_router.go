package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/equipment/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewEquipmentRouter(handler interfaces.EquipmentHandlerInterface) http.Handler {
	r := chi.NewRouter()

	// Equipment CRUD endpoints
	r.Post("/", handler.CreateEquipment)       // POST /equipment
	r.Get("/", handler.ListEquipment)          // GET /equipment
	r.Get("/{id}", handler.GetEquipment)       // GET /equipment/{id}
	r.Put("/{id}", handler.UpdateEquipment)    // PUT /equipment/{id}
	r.Delete("/{id}", handler.DeleteEquipment) // DELETE /equipment/{id}

	return r
}
