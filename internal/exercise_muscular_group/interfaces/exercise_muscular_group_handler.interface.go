package interfaces

import "net/http"

// ExerciseMuscularGroupHandler defines handler contract for join table endpoints

type ExerciseMuscularGroupHandler interface {
	CreateLink(w http.ResponseWriter, r *http.Request)
	DeleteLink(w http.ResponseWriter, r *http.Request)
	GetLinkByID(w http.ResponseWriter, r *http.Request)
	GetLinksByExerciseID(w http.ResponseWriter, r *http.Request)
	GetLinksByMuscularGroupID(w http.ResponseWriter, r *http.Request)
}
