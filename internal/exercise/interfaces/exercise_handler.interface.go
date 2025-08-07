package interfaces

import "net/http"

type ExerciseHandler interface {
	// Exercise management
	CreateExercise(w http.ResponseWriter, r *http.Request)
	GetExerciseByID(w http.ResponseWriter, r *http.Request)
	GetExerciseByMuscularGroup(w http.ResponseWriter, r *http.Request)
	GetExerciseByEquipment(w http.ResponseWriter, r *http.Request)
	GetAllExercises(w http.ResponseWriter, r *http.Request)
	UpdateExercise(w http.ResponseWriter, r *http.Request)
	DeleteExercise(w http.ResponseWriter, r *http.Request)
	GetExercisesByFilters(w http.ResponseWriter, r *http.Request)
}
