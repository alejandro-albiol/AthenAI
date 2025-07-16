package interfaces

import "net/http"

type ExerciseHandler interface {
	// Exercise management
	CreateExercise(w http.ResponseWriter, r *http.Request)
	GetExerciseByID(w http.ResponseWriter, r *http.Request)
	GetAllExercises(w http.ResponseWriter, r *http.Request)
	UpdateExercise(w http.ResponseWriter, r *http.Request)
	DeleteExercise(w http.ResponseWriter, r *http.Request)
	GetExercisesByMuscularGroup(w http.ResponseWriter, r *http.Request)
	GetExercisesByEquipment(w http.ResponseWriter, r *http.Request)
}
