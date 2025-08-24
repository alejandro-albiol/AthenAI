package router

import (
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/handler"
	"github.com/go-chi/chi/v5"
)

func NewCustomMemberWorkoutRouter(h *handler.CustomMemberWorkoutHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/custom-member-workout", h.Create)
	r.Get("/custom-member-workout/{id}", h.GetByID)
	r.Get("/custom-member-workout/member/{memberID}", h.ListByMemberID)
	r.Put("/custom-member-workout/{id}", h.Update)
	r.Delete("/custom-member-workout/{id}", h.Delete)
	return r
}
