package router

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
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
		id := chi.URLParam(r, "id")
		handler.GetUserByID(w, r, id)
	})

	r.Get("/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		handler.GetUserByUsername(w, r, username)
	})

	r.Get("/email/{email}", func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		handler.GetUserByEmail(w, r, email)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var userDTO dto.UserUpdateDTO
		if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handler.UpdateUser(w, r, id, userDTO)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.DeleteUser(w, r, id)
	})

	r.Post("/{id}/verify", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler.VerifyUser(w, r, id)
	})

	r.Post("/{id}/active", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var activeReq struct {
			Active bool `json:"active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&activeReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handler.SetUserActive(w, r, id, activeReq.Active)
	})

	return r
}
