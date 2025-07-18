package interfaces

import "net/http"

type ExerciseHandler interface {
	// Exercise management
	CreateExercise(w http.ResponseWriter, r *http.Request)
	GetExerciseByID(w http.ResponseWriter, r *http.Request, id string)
	GetExerciseByMuscularGroup(w http.ResponseWriter, r *http.Request, groups []string)
	GetExerciseByEquipment(w http.ResponseWriter, r *http.Request, equipment []string)
	GetAllExercises(w http.ResponseWriter, r *http.Request)
	UpdateExercise(w http.ResponseWriter, r *http.Request, id string)
	DeleteExercise(w http.ResponseWriter, r *http.Request, id string)
	GetExercisesByFilters(w http.ResponseWriter, r *http.Request, groups []string, equipment []string)
}
