package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/router"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/service"
)

func NewCustomMemberWorkoutModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomMemberWorkoutRepository(db)
	service := service.NewCustomMemberWorkoutService(repo)
	handler := handler.NewCustomMemberWorkoutHandler(service)
	return router.NewCustomMemberWorkoutRouter(handler)
}
