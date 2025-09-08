package interfaces

import (
	"net/http"
)

// CustomExerciseHandler defines HTTP handler interface for custom exercises

type CustomExerciseHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
