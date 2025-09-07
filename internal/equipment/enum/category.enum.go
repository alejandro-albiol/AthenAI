package enum

type EquipmentCategory string

const (
	EquipmentCategoryFreeWeights EquipmentCategory = "free_weights"
	EquipmentCategoryMachines    EquipmentCategory = "machines"
	EquipmentCategoryCardio      EquipmentCategory = "cardio"
	EquipmentCategoryAccessories EquipmentCategory = "accessories"
	EquipmentCategoryBodyweight  EquipmentCategory = "bodyweight"
)

func (e EquipmentCategory) IsValid() bool {
	switch e {
	case EquipmentCategoryFreeWeights, EquipmentCategoryMachines, EquipmentCategoryCardio, EquipmentCategoryAccessories, EquipmentCategoryBodyweight:
		return true
	}
	return false
}
