package dto

type WorkoutGeneratorRequest struct {
	UserID           string   `json:"user_id"`
	Tags             []string `json:"tags"`
	TrainingPhase    string   `json:"training_phase"`
	Motivation       string   `json:"motivation"`
	SpecialSituation string   `json:"special_situation"`
	TemplateName     string   `json:"template_name"`
}
