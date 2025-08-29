package interfaces

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
)

type CustomExerciseMuscularGroupHandler interface {
	CreateLink(r *http.Request) error
	DeleteLink(r *http.Request) error
	RemoveAllLinksForExercise(r *http.Request) error
	GetLinkByID(r *http.Request) (dto.CustomExerciseMuscularGroup, error)
	GetLinksByCustomExerciseID(r *http.Request) ([]dto.CustomExerciseMuscularGroup, error)
	GetLinksByMuscularGroupID(r *http.Request) ([]dto.CustomExerciseMuscularGroup, error)
}
