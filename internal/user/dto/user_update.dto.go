package dto

type UserUpdateDTO struct {
	Username         string  `json:"username,omitempty"`
	Email            string  `json:"email,omitempty"`
	Description      *string `json:"description,omitempty"`
	TrainingPhase    *string `json:"training_phase,omitempty" validate:"omitempty,oneof=weight_loss muscle_gain cardio_improve maintenance"`
	Motivation       *string `json:"motivation,omitempty" validate:"omitempty,oneof=medical_recommendation self_improvement competition rehabilitation wellbeing"`
	SpecialSituation *string `json:"special_situation,omitempty" validate:"omitempty,oneof=pregnancy post_partum injury_recovery chronic_condition elderly_population physical_limitation none"`
}
