package api

import (
	"database/sql"

	"net/http"

	authmodule "github.com/alejandro-albiol/athenai/internal/auth/module"
	gymmodule "github.com/alejandro-albiol/athenai/internal/gym/module"
	usermodule "github.com/alejandro-albiol/athenai/internal/user/module"

	"github.com/go-chi/chi/v5"

	equipmentmodule "github.com/alejandro-albiol/athenai/internal/equipment/module"
	exercisemodule "github.com/alejandro-albiol/athenai/internal/exercise/module"
	exerciseequipmentmodule "github.com/alejandro-albiol/athenai/internal/exercise_equipment/module"
	exercisemuscgroupmodule "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/module"
	musculargroupmodule "github.com/alejandro-albiol/athenai/internal/muscular_group/module"
	templateblockmodule "github.com/alejandro-albiol/athenai/internal/template_block/module"
	workouttemplatemodule "github.com/alejandro-albiol/athenai/internal/workout_template/module"
	// adminmodule "github.com/alejandro-albiol/athenai/internal/admin/module"
	// customexercisemuscgroupmodule "github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/module"
)

func NewAPIRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// Auth
	auth := authmodule.NewAuthModule(db)
	r.Mount("/auth", auth.Router)

	// Gym
	r.Mount("/gym", gymmodule.NewGymModule(db))

	// User
	r.Mount("/user", usermodule.NewUserModule(db))

	// Equipment
	r.Mount("/equipment", equipmentmodule.NewEquipmentModule(db))

	// Exercise
	r.Mount("/exercise", exercisemodule.NewExerciseModule(db))

	// Workout Generator
	// NOTE: This module may require service dependencies from other modules.
	// If so, wire them up here before passing to NewWorkoutGeneratorModule.
	// For now, assuming direct construction is possible:
	// r.Mount("/workout-generator", workoutgeneratormodule.NewWorkoutGeneratorModule(...))

	// Workout Template
	r.Mount("/workout-template", workouttemplatemodule.NewWorkoutTemplateModule(db))

	// Template Block
	r.Mount("/template-block", templateblockmodule.NewTemplateBlockModule(db))

	// Muscular Group
	r.Mount("/muscular-group", musculargroupmodule.NewMuscularGroupModule(db))

	// Exercise Equipment
	r.Mount("/exercise-equipment", exerciseequipmentmodule.NewExerciseEquipmentModule(db))

	// Exercise Muscular Group
	r.Mount("/exercise-muscular-group", exercisemuscgroupmodule.NewExerciseMuscularGroupModule(db))

	// Uncomment when implemented:
	// r.Mount("/admin", adminmodule.NewAdminModule(db))
	// r.Mount("/custom-exercise-muscular-group", customexercisemuscgroupmodule.NewCustomExerciseMuscularGroupModule(db))

	return r
}
