package enum

type TrainingPhase string

const (
	WeightLoss    TrainingPhase = "weight_loss"
	MuscleGain    TrainingPhase = "muscle_gain"
	CardioImprove TrainingPhase = "cardio_improve"
	Maintenance   TrainingPhase = "maintenance"
)

func (tp TrainingPhase) IsValid() bool {
	switch tp {
	case WeightLoss, MuscleGain, CardioImprove, Maintenance:
		return true
	}
	return false
}
