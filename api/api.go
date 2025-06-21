package api

import (
	"github.com/alejandro-albiol/athenai/internal/databases/interfaces"
	usersmodule "github.com/alejandro-albiol/athenai/internal/users/modules"
	"github.com/go-chi/chi/v5"
)

func NewAPIModule(db interfaces.DBService) *chi.Mux {
    r := chi.NewRouter()
    r.Mount("/users", usersmodule.NewUserModule(db))
    //r.Mount("/exercises", module.NewExercisesModule(db))
    // Add more modules as needed
    return r
}