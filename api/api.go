package api

import (
	"database/sql"

	usermodule "github.com/alejandro-albiol/athenai/internal/user/module"
	"github.com/go-chi/chi/v5"
)

func NewAPIModule(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/users", usermodule.NewUserModule(db))
	// r.Mount("/exercises", module.NewExercisesModule(db))
	// Add more modules as needed
	return r
}
