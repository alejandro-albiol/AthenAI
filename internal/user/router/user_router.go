package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewUsersRouter(handler interfaces.UserHandler) http.Handler {
	r := chi.NewRouter()

	// Auth middleware is applied globally at the API level
	// All routes here already have authenticated user context with gym ID

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handler.RegisterUser(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAllUsers(w, r)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetUserByID(w, r)
	})

	r.Get("/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetUserByUsername(w, r)
	})

	r.Get("/email/{email}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetUserByEmail(w, r)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateUser(w, r)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteUser(w, r)
	})

	r.Post("/{id}/verify", func(w http.ResponseWriter, r *http.Request) {
		handler.VerifyUser(w, r)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteUser(w, r)
	})

	r.Post("/{id}/verify", func(w http.ResponseWriter, r *http.Request) {
		handler.VerifyUser(w, r)
	})

	r.Post("/{id}/active", func(w http.ResponseWriter, r *http.Request) {
		handler.SetUserActive(w, r)
	})

	return r
}
