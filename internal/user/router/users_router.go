package router

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/handler"
	"github.com/go-chi/chi/v5"
)

func NewUsersRouter(handler *handler.UsersHandler) http.Handler {
	r := chi.NewRouter()

	getGymID := func(r *http.Request) string {
		return r.Header.Get("X-Gym-ID")
	}

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		gymID := getGymID(r)
		handler.RegisterUser(w, r, gymID)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		gymID := getGymID(r)
		handler.GetAllUsers(w, gymID)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		gymID := getGymID(r)
		id := chi.URLParam(r, "id")
		handler.GetUserByID(w, gymID, id)
	})

	r.Get("/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		gymID := getGymID(r)
		username := chi.URLParam(r, "username")
		handler.GetUserByUsername(w, gymID, username)
	})

	r.Get("/email/{email}", func(w http.ResponseWriter, r *http.Request) {
		gymID := getGymID(r)
		email := chi.URLParam(r, "email")
		handler.GetUserByEmail(w, gymID, email)
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		gymID := getGymID(r)
		id := chi.URLParam(r, "id")
		var userDTO dto.UserUpdateDTO
		if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handler.UpdateUser(w, gymID, id, userDTO)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		gymID := getGymID(r)
		id := chi.URLParam(r, "id")
		handler.DeleteUser(w, gymID, id)
	})

	return r
}
