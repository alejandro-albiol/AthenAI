package service

import (
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	mgdto "github.com/alejandro-albiol/athenai/internal/muscular_group/dto"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	CreateLinkFn                func(gymID string, req *dto.CustomExerciseMuscularGroupCreationDTO) (*string, error)
	DeleteLinkFn                func(gymID, id string) error
	RemoveAllLinksForExerciseFn func(gymID, customExerciseID string) error
	FindByIDFn                  func(gymID, id string) (*dto.CustomExerciseMuscularGroup, error)
	FindByCustomExerciseIDFn    func(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error)
	FindByMuscularGroupIDFn     func(gymID, muscularGroupID string) ([]*dto.CustomExerciseMuscularGroup, error)
}

func (m *mockRepo) CreateLink(gymID string, req *dto.CustomExerciseMuscularGroupCreationDTO) (*string, error) {
	return m.CreateLinkFn(gymID, req)
}
func (m *mockRepo) DeleteLink(gymID, id string) error {
	return m.DeleteLinkFn(gymID, id)
}
func (m *mockRepo) RemoveAllLinksForExercise(gymID, customExerciseID string) error {
	return m.RemoveAllLinksForExerciseFn(gymID, customExerciseID)
}
func (m *mockRepo) FindByID(gymID, id string) (*dto.CustomExerciseMuscularGroup, error) {
	return m.FindByIDFn(gymID, id)
}
func (m *mockRepo) FindByCustomExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error) {
	return m.FindByCustomExerciseIDFn(gymID, customExerciseID)
}
func (m *mockRepo) FindByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.CustomExerciseMuscularGroup, error) {
	return m.FindByMuscularGroupIDFn(gymID, muscularGroupID)
}

type mockPublicMuscularGroupRepo struct{}

func (m *mockPublicMuscularGroupRepo) CreateMuscularGroup(req *mgdto.CreateMuscularGroupDTO) (*string, error) {
	return nil, nil
}
func (m *mockPublicMuscularGroupRepo) GetAllMuscularGroups() ([]*mgdto.MuscularGroupResponseDTO, error) {
	return nil, nil
}
func (m *mockPublicMuscularGroupRepo) GetMuscularGroupByID(id string) (*mgdto.MuscularGroupResponseDTO, error) {
	return &mgdto.MuscularGroupResponseDTO{ID: id, IsActive: true}, nil
}
func (m *mockPublicMuscularGroupRepo) GetMuscularGroupByName(name string) (*mgdto.MuscularGroupResponseDTO, error) {
	return nil, nil
}
func (m *mockPublicMuscularGroupRepo) UpdateMuscularGroup(id string, mg *mgdto.UpdateMuscularGroupDTO) (*mgdto.MuscularGroupResponseDTO, error) {
	return nil, nil
}
func (m *mockPublicMuscularGroupRepo) DeleteMuscularGroup(id string) error { return nil }
func TestCreateLink_Success(t *testing.T) {
	repo := &mockRepo{
		CreateLinkFn: func(gymID string, req *dto.CustomExerciseMuscularGroupCreationDTO) (*string, error) {
			id := "123"
			return &id, nil
		},
	}
	svc := &CustomExerciseMuscularGroupService{
		repository:              repo,
		publicMuscularGroupRepo: &mockPublicMuscularGroupRepo{},
	}
	id, err := svc.CreateLink("tenant1", &dto.CustomExerciseMuscularGroupCreationDTO{CustomExerciseID: "ex1", MuscularGroupID: "mg1"})
	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestDeleteLink_Success(t *testing.T) {
	repo := &mockRepo{
		DeleteLinkFn: func(gymID, id string) error { return nil },
	}
	svc := &CustomExerciseMuscularGroupService{repository: repo, publicMuscularGroupRepo: &mockPublicMuscularGroupRepo{}}
	err := svc.DeleteLink("tenant1", "id1")
	assert.NoError(t, err)
}

func TestRemoveAllLinksForExercise_Success(t *testing.T) {
	repo := &mockRepo{
		RemoveAllLinksForExerciseFn: func(gymID, customExerciseID string) error { return nil },
	}
	svc := &CustomExerciseMuscularGroupService{repository: repo, publicMuscularGroupRepo: &mockPublicMuscularGroupRepo{}}
	err := svc.RemoveAllLinksForExercise("tenant1", "ex1")
	assert.NoError(t, err)
}

func TestGetLinkByID_Success(t *testing.T) {
	repo := &mockRepo{
		FindByIDFn: func(gymID, id string) (*dto.CustomExerciseMuscularGroup, error) {
			return &dto.CustomExerciseMuscularGroup{ID: id, CustomExerciseID: "ex1", MuscularGroupID: "mg1"}, nil
		},
	}
	svc := &CustomExerciseMuscularGroupService{repository: repo, publicMuscularGroupRepo: &mockPublicMuscularGroupRepo{}}
	link, err := svc.GetLinkByID("tenant1", "id1")
	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, "id1", link.ID)
}

func TestGetLinksByCustomExerciseID_Success(t *testing.T) {
	repo := &mockRepo{
		FindByCustomExerciseIDFn: func(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error) {
			return []*dto.CustomExerciseMuscularGroup{{ID: "id1", CustomExerciseID: customExerciseID, MuscularGroupID: "mg1"}}, nil
		},
	}
	svc := &CustomExerciseMuscularGroupService{repository: repo, publicMuscularGroupRepo: &mockPublicMuscularGroupRepo{}}
	links, err := svc.GetLinksByCustomExerciseID("tenant1", "ex1")
	assert.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, "ex1", links[0].CustomExerciseID)
}

func TestGetLinksByMuscularGroupID_Success(t *testing.T) {
	repo := &mockRepo{
		FindByMuscularGroupIDFn: func(gymID, muscularGroupID string) ([]*dto.CustomExerciseMuscularGroup, error) {
			return []*dto.CustomExerciseMuscularGroup{{ID: "id1", CustomExerciseID: "ex1", MuscularGroupID: muscularGroupID}}, nil
		},
	}
	svc := &CustomExerciseMuscularGroupService{repository: repo, publicMuscularGroupRepo: &mockPublicMuscularGroupRepo{}}
	links, err := svc.GetLinksByMuscularGroupID("tenant1", "mg1")
	assert.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, "mg1", links[0].MuscularGroupID)
}

// ...additional tests for DeleteLink, RemoveAllLinksForExercise, GetLinkByID, GetLinksByExerciseID, GetLinksByMuscularGroupID...
