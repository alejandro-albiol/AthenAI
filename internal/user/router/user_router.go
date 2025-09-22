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

	// Platform admin routes - allow specifying gym context
	r.Route("/gym/{gymId}", func(r chi.Router) {
		r.Get("/", handler.GetUsersByGymID)                // GET /user/gym/{gymId} - Platform admin only
		r.Post("/", handler.RegisterUserInGym)             // POST /user/gym/{gymId} - Platform admin only
		r.Get("/{id}", handler.GetUserByIDInGym)           // GET /user/gym/{gymId}/{id} - Platform admin only
		r.Put("/{id}", handler.UpdateUserInGym)            // PUT /user/gym/{gymId}/{id} - Platform admin only
		r.Delete("/{id}", handler.DeleteUserInGym)         // DELETE /user/gym/{gymId}/{id} - Platform admin only
		r.Post("/{id}/verify", handler.VerifyUserInGym)    // POST /user/gym/{gymId}/{id}/verify - Platform admin only
		r.Post("/{id}/active", handler.SetUserActiveInGym) // POST /user/gym/{gymId}/{id}/active - Platform admin only
	})

	// Regular user CRUD endpoints (gym context from JWT)
	r.Post("/", handler.RegisterUser)                        // POST /user
	r.Get("/", handler.GetAllUsers)                          // GET /user
	r.Get("/{id}", handler.GetUserByID)                      // GET /user/{id}
	r.Get("/username/{username}", handler.GetUserByUsername) // GET /user/username/{username}
	r.Get("/email/{email}", handler.GetUserByEmail)          // GET /user/email/{email}
	r.Put("/{id}", handler.UpdateUser)                       // PUT /user/{id}
	r.Delete("/{id}", handler.DeleteUser)                    // DELETE /user/{id}
	r.Post("/{id}/verify", handler.VerifyUser)               // POST /user/{id}/verify
	r.Post("/{id}/active", handler.SetUserActive)            // POST /user/{id}/active

	return r
}
