package enum

type SpecialSituation string

const (
	Pregnancy          SpecialSituation = "pregnancy"
	PostPartum         SpecialSituation = "post_partum"
	InjuryRecovery     SpecialSituation = "injury_recovery"
	ChronicCondition   SpecialSituation = "chronic_condition"
	ElderlyPopulation  SpecialSituation = "elderly_population"
	PhysicalLimitation SpecialSituation = "physical_limitation"
	None               SpecialSituation = "none"
)

func (s SpecialSituation) IsValid() bool {
	switch s {
	case Pregnancy, PostPartum, InjuryRecovery, ChronicCondition, ElderlyPopulation, PhysicalLimitation, None:
		return true
	}
	return false
}
