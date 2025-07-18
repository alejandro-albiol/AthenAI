package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/handler"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/repository"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/router"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/service"
)

func NewExerciseEquipmentModule(db *sql.DB) http.Handler {
	repo := repository.NewExerciseEquipmentRepository(db)
	service := service.NewExerciseEquipmentService(repo)
	handler := handler.NewExerciseEquipmentHandler(service)
	return router.NewExerciseEquipmentRouter(handler)
}
