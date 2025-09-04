package interfaces

import "github.com/alejandro-albiol/athenai/internal/muscular_group/dto"

type MuscularGroupRepository interface {
	CreateMuscularGroup(mg *dto.CreateMuscularGroupDTO) (string, error)
	GetAllMuscularGroups() ([]*dto.MuscularGroupResponseDTO, error)
	GetMuscularGroupByID(id string) (*dto.MuscularGroupResponseDTO, error)
	GetMuscularGroupByName(name string) (*dto.MuscularGroupResponseDTO, error)
	UpdateMuscularGroup(id string, mg *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error)
	DeleteMuscularGroup(id string) error
}
