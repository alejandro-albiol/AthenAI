package interfaces

import "github.com/alejandro-albiol/athenai/internal/muscular_group/dto"

type MuscularGroupService interface {
	CreateMuscularGroup(mg dto.MuscularGroup) (string, error)
	GetAllMuscularGroups() ([]dto.MuscularGroup, error)
	UpdateMuscularGroup(id string, mg dto.MuscularGroup) (dto.MuscularGroup, error)
	DeleteMuscularGroup(id string) error
	FindByID(id string) (dto.MuscularGroup, error)
}