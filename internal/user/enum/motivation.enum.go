package enum

type Motivation string

const (
	MedicalRecommendation Motivation = "medical_recommendation"
	SelfImprovement       Motivation = "self_improvement"
	Competition           Motivation = "competition"
	Rehabilitation        Motivation = "rehabilitation"
	Wellbeing             Motivation = "wellbeing"
)

func (m Motivation) IsValid() bool {
	switch m {
	case MedicalRecommendation, SelfImprovement, Competition, Rehabilitation, Wellbeing:
		return true
	}
	return false
}
