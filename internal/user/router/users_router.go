package router

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func NewUsersRouter(handler interfaces.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handler.RegisterUser(w, r, middleware.GetGymID(r))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAllUsers(w, middleware.GetGymID(r))
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.GetUserByID(w, middleware.GetGymID(r), id)
	})

	r.Get("/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		handler.GetUserByUsername(w, middleware.GetGymID(r), username)
	})

	r.Get("/email/{email}", func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		handler.GetUserByEmail(w, middleware.GetGymID(r), email)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var userDTO dto.UserUpdateDTO
		if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handler.UpdateUser(w, middleware.GetGymID(r), id, userDTO)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.DeleteUser(w, middleware.GetGymID(r), id)
	})

	r.Post("/{id}/verify", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.VerifyUser(w, middleware.GetGymID(r), id)
	})

	r.Post("/{id}/active", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var active bool
		if err := json.NewDecoder(r.Body).Decode(&active); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handler.SetUserActive(w, middleware.GetGymID(r), id, active)
	})

	return r
}
