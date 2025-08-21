package module

import (
	"database/sql"
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/exercise/handler"
	"github.com/alejandro-albiol/athenai/internal/exercise/repository"
	"github.com/alejandro-albiol/athenai/internal/exercise/router"
	"github.com/alejandro-albiol/athenai/internal/exercise/service"
	exerciseMuscularGroupRepository "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/repository"
	exerciseEquipmentRepository "github.com/alejandro-albiol/athenai/internal/exercise_equipment/repository"
	exerciseMuscularGroupService "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/service"
	exerciseEquipmentService "github.com/alejandro-albiol/athenai/internal/exercise_equipment/service"
)

func NewExerciseModule(db *sql.DB) http.Handler {
	repo := repository.NewExerciseRepository(db)
	exerciseMuscularGroupRepository := exerciseMuscularGroupRepository.NewExerciseMuscularGroupRepository(db)
	exerciseEquipmentRepository := exerciseEquipmentRepository.NewExerciseEquipmentRepository(db)
	exerciseEquipmentService := exerciseEquipmentService.NewExerciseEquipmentService(exerciseEquipmentRepository)
	exerciseMuscularGroupService := exerciseMuscularGroupService.NewExerciseMuscularGroupService(exerciseMuscularGroupRepository)
	service := service.NewExerciseService(repo, exerciseEquipmentService, exerciseMuscularGroupService)
	handler := handler.NewExerciseHandler(service)
	return router.NewExerciseRouter(handler)
}
