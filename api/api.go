package api

import (
	"database/sql"

	authmodule "github.com/alejandro-albiol/athenai/internal/auth/module"
	gymmodule "github.com/alejandro-albiol/athenai/internal/gym/module"
	usermodule "github.com/alejandro-albiol/athenai/internal/user/module"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func NewAPIModule(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// Create auth module (contains both service and router)
	authMod := authmodule.NewAuthModule(db)

	// Mount auth routes first (no auth middleware for login/register)
	r.Mount("/auth", authMod.Router)

	// Apply global auth middleware to all other routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authMod.Service))

		// All these routes require authentication
		r.Mount("/user", usermodule.NewUserModule(db))
		r.Mount("/gym", gymmodule.NewGymModule(db))
		// r.Mount("/exercises", module.NewExercisesModule(db))
		// Add more modules as needed
	})

	return r
}
