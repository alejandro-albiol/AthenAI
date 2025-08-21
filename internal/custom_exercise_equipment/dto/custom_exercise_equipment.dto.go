package dto

type CustomExerciseEquipment struct {
	ID               string `json:"id"`
	CustomExerciseID string `json:"custom_exercise_id"`
	EquipmentID      string `json:"equipment_id"`
}
