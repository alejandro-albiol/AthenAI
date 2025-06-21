package router

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/users/handlers"
	"github.com/alejandro-albiol/athenai/internal/users/interfaces"
	"github.com/go-chi/chi/v5"
)

func NewUsersRouter(handler *handlers.UsersHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/", handler.RegisterUser)
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.GetUserByID(w, id)
	})
	r.Get("/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		handler.GetUserByUsername(w, username)
	})
	r.Get("/email/{email}", func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		handler.GetUserByEmail(w, email)
	})
	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		var user interfaces.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.ID = chi.URLParam(r, "id")
		handler.UpdateUser(w, user)
	})
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.DeleteUser(w, id)
	})

	return r
}
