package module

import (
	"database/sql"
	"net/http"

	custom_equipment_repository "github.com/alejandro-albiol/athenai/internal/custom_equipment/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/handler"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/repository"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/router"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/service"
	equipment_repository "github.com/alejandro-albiol/athenai/internal/equipment/repository"
)

type CustomExerciseEquipmentModule struct {
	Handler *handler.CustomExerciseEquipmentHandler
	Router  http.Handler
}

func NewCustomExerciseEquipmentModule(db *sql.DB) http.Handler {
	repo := repository.NewCustomExerciseEquipmentRepository(db)
	customEquipmentRepo := custom_equipment_repository.NewCustomEquipmentRepository(db)
	publicEquipmentRepo := equipment_repository.NewEquipmentRepository(db)
	svc := service.NewCustomExerciseEquipmentService(repo, customEquipmentRepo, publicEquipmentRepo)
	handler := handler.NewCustomExerciseEquipmentHandler(svc)
	router := router.NewCustomExerciseEquipmentRouter(handler)
	return router
}
