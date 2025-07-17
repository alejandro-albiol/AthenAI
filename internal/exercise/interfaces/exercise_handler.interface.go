package interfaces

import "net/http"

type ExerciseHandler interface {
	// Exercise management
	CreateExercise(w http.ResponseWriter, r *http.Request)
	GetExerciseByID(w http.ResponseWriter, r *http.Request, id string)
	GetAllExercises(w http.ResponseWriter, r *http.Request)
	UpdateExercise(w http.ResponseWriter, r *http.Request, id string)
	DeleteExercise(w http.ResponseWriter, r *http.Request, id string)
	GetExercisesByMuscularGroup(w http.ResponseWriter, r *http.Request, group string)
	GetExercisesByEquipment(w http.ResponseWriter, r *http.Request, equipment string)
}
