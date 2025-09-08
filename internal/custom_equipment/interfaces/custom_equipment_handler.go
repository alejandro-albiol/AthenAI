package interfaces

import (
	"net/http"
)

// CustomEquipmentHandler defines HTTP handler interface for custom equipment

type CustomEquipmentHandler interface {
	CreateCustomEquipment(w http.ResponseWriter, r *http.Request)
	GetCustomEquipmentByID(w http.ResponseWriter, r *http.Request)
	GetCustomEquipmentByName(w http.ResponseWriter, r *http.Request)
	ListCustomEquipment(w http.ResponseWriter, r *http.Request)
	UpdateCustomEquipment(w http.ResponseWriter, r *http.Request)
	DeleteCustomEquipment(w http.ResponseWriter, r *http.Request)
}
