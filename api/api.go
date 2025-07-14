package api

import (
	"database/sql"

	authmodule "github.com/alejandro-albiol/athenai/internal/auth/module"
	gymmodule "github.com/alejandro-albiol/athenai/internal/gym/module"
	usermodule "github.com/alejandro-albiol/athenai/internal/user/module"
	"github.com/go-chi/chi/v5"
)

func NewAPIModule(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	// r.Use(middleware.RequireGymID)

	// Mount routes
	r.Mount("/auth", authmodule.NewAuthModule(db))
	r.Mount("/user", usermodule.NewUserModule(db))
	r.Mount("/gym", gymmodule.NewGymModule(db))
	// r.Mount("/exercises", module.NewExercisesModule(db))
	// Add more modules as needed

	return r
}
