package interfaces

import (
	"net/http"
)

type CustomExerciseMuscularGroupHandler interface {
	CreateLink(w http.ResponseWriter, r *http.Request)
	DeleteLink(w http.ResponseWriter, r *http.Request)
	GetLinkByID(w http.ResponseWriter, r *http.Request)
	GetLinksByExerciseID(w http.ResponseWriter, r *http.Request)
	GetLinksByMuscularGroupID(w http.ResponseWriter, r *http.Request)
	RemoveAllLinksForExercise(w http.ResponseWriter, r *http.Request)
}
