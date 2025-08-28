package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/router"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/service"
)

type CustomExerciseEquipmentModule struct {
	Handler *handler.CustomExerciseEquipmentHandler
	Router  http.Handler
}

func NewCustomExerciseEquipmentModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomExerciseEquipmentRepository(db)
	svc := service.NewCustomExerciseEquipmentService(repo)
	handler := handler.NewCustomExerciseEquipmentHandler(svc)
	router := router.NewCustomExerciseEquipmentRouter(handler)

	return router
}
