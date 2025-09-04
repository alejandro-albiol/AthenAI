package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/muscular_group/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewMuscularGroupRouter(handler interfaces.MuscularGroupHandlerInterface) http.Handler {
	r := chi.NewRouter()

	r.Post("/", handler.CreateMuscularGroup)       // POST /muscular-groups
	r.Get("/", handler.ListMuscularGroups)         // GET /muscular-groups
	r.Get("/{id}", handler.GetMuscularGroup)       // GET /muscular-groups/{id}
	r.Put("/{id}", handler.UpdateMuscularGroup)    // PUT /muscular-groups/{id}
	r.Delete("/{id}", handler.DeleteMuscularGroup) // DELETE /muscular-groups/{id}

	return r
}
