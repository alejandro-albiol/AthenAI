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

	// User CRUD endpoints
	r.Post("/", handler.RegisterUser)                        // POST /users
	r.Get("/", handler.GetAllUsers)                          // GET /users
	r.Get("/{id}", handler.GetUserByID)                      // GET /users/{id}
	r.Get("/username/{username}", handler.GetUserByUsername) // GET /users/username/{username}
	r.Get("/email/{email}", handler.GetUserByEmail)          // GET /users/email/{email}
	r.Put("/{id}", handler.UpdateUser)                       // PUT /users/{id}
	r.Delete("/{id}", handler.DeleteUser)                    // DELETE /users/{id}
	r.Post("/{id}/verify", handler.VerifyUser)               // POST /users/{id}/verify
	r.Post("/{id}/active", handler.SetUserActive)            // POST /users/{id}/active

	return r
}
