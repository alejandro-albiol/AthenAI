package interfaces

import "github.com/alejandro-albiol/athenai/internal/muscular_group/dto"

type MuscularGroupService interface {
	CreateMuscularGroup(mg *dto.CreateMuscularGroupDTO) (*string, error)
	GetMuscularGroupByID(id string) (*dto.MuscularGroupResponseDTO, error)
	GetAllMuscularGroups() ([]*dto.MuscularGroupResponseDTO, error)
	UpdateMuscularGroup(id string, mg *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error)
	DeleteMuscularGroup(id string) error
}
