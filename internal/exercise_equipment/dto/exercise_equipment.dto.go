package dto

type ExerciseEquipment struct {
	ID          string `json:"id"`
	ExerciseID  string `json:"exercise_id"`
	EquipmentID string `json:"equipment_id"`
}
