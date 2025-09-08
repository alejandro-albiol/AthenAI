package interfaces

import "net/http"

type CustomExerciseEquipmentHandler interface {
	CreateLink(w http.ResponseWriter, r *http.Request)
	DeleteLink(w http.ResponseWriter, r *http.Request)
	RemoveAllLinksForExercise(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	FindByCustomExerciseID(w http.ResponseWriter, r *http.Request)
	FindByEquipmentID(w http.ResponseWriter, r *http.Request)
}
