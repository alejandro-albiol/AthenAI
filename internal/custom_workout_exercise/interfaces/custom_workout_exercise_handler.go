package interfaces

import (
	"net/http"
)

type CustomWorkoutExerciseHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	ListByWorkoutInstanceID(w http.ResponseWriter, r *http.Request)
	ListByEquipmentID(w http.ResponseWriter, r *http.Request)
	ListByMuscularGroupID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
