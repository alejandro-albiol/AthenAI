package api

import (
	"database/sql"

	"net/http"

	authmodule "github.com/alejandro-albiol/athenai/internal/auth/module"
	gymmodule "github.com/alejandro-albiol/athenai/internal/gym/module"
	usermodule "github.com/alejandro-albiol/athenai/internal/user/module"
	"github.com/alejandro-albiol/athenai/pkg/middleware"

	"github.com/go-chi/chi/v5"

	equipmentmodule "github.com/alejandro-albiol/athenai/internal/equipment/module"
	exercisemodule "github.com/alejandro-albiol/athenai/internal/exercise/module"
	exerciseequipmentmodule "github.com/alejandro-albiol/athenai/internal/exercise_equipment/module"
	exercisemuscgroupmodule "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/module"
	musculargroupmodule "github.com/alejandro-albiol/athenai/internal/muscular_group/module"
	templateblockmodule "github.com/alejandro-albiol/athenai/internal/template_block/module"
	workouttemplatemodule "github.com/alejandro-albiol/athenai/internal/workout_template/module"
	customequipmentmodule "github.com/alejandro-albiol/athenai/internal/custom_equipment/module"
	customexercisemodule "github.com/alejandro-albiol/athenai/internal/custom_exercise/module"
	// workoutgeneratormodule "github.com/alejandro-albiol/athenai/internal/workout_generator/module"
	// adminmodule "github.com/alejandro-albiol/athenai/internal/admin/module"
	// customexercisemuscgroupmodule "github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/module"
)

func NewAPIRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// Mount public auth routes
	auth := authmodule.NewAuthModule(db)
	r.Mount("/auth", auth.Router)

	// Protected routes subrouter
	protected := chi.NewRouter()
	protected.Use(middleware.AuthMiddleware(auth.Service))
	protected.Mount("/gym", gymmodule.NewGymModule(db))
	protected.Mount("/user", usermodule.NewUserModule(db))
	protected.Mount("/equipment", equipmentmodule.NewEquipmentModule(db))
	protected.Mount("/exercise", exercisemodule.NewExerciseModule(db))
	// protected.Mount("/workout-generator", workoutgeneratormodule.NewWorkoutGeneratorModule(...))
	protected.Mount("/workout-template", workouttemplatemodule.NewWorkoutTemplateModule(db))
	protected.Mount("/template-block", templateblockmodule.NewTemplateBlockModule(db))
	protected.Mount("/muscular-group", musculargroupmodule.NewMuscularGroupModule(db))
	protected.Mount("/exercise-equipment", exerciseequipmentmodule.NewExerciseEquipmentModule(db))
	protected.Mount("/exercise-muscular-group", exercisemuscgroupmodule.NewExerciseMuscularGroupModule(db))
	// Uncomment when implemented:
	// protected.Mount("/admin", adminmodule.NewAdminModule(db))
	// protected.Mount("/custom-exercise-muscular-group", customexercisemuscgroupmodule.NewCustomExerciseMuscularGroupModule(db))
	protected.Mount("/custom-equipment", customequipmentmodule.NewCustomEquipmentModule(db))
	protected.Mount("/custom-exercise", customexercisemodule.NewCustomExerciseModule(db))

	r.Mount("/", protected)
	return r
}
