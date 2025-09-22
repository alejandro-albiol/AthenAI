package dto

import userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"

type UserCreationDTO struct {
	Username         string                 `json:"username"`
	Email            string                 `json:"email"`
	Password         string                 `json:"password"`
	Role             userrole_enum.UserRole `json:"role,omitempty"`
	Description      *string                `json:"description,omitempty"`
	TrainingPhase    string                 `json:"training_phase,omitempty" validate:"omitempty,oneof=weight_loss muscle_gain cardio_improve maintenance"`
	Motivation       string                 `json:"motivation,omitempty" validate:"omitempty,oneof=medical_recommendation self_improvement competition rehabilitation wellbeing"`
	SpecialSituation string                 `json:"special_situation,omitempty" validate:"omitempty,oneof=pregnancy post_partum injury_recovery chronic_condition elderly_population physical_limitation none"`
}
