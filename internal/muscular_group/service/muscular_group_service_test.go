package service_test

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	groups       []*dto.MuscularGroupResponseDTO
	createErr    error
	getAllErr    error
	getByIDErr   error
	updateErr    error
	deleteErr    error
	getByNameErr error
}

func (m *mockRepository) CreateMuscularGroup(mg *dto.CreateMuscularGroupDTO) (*string, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	id := "new-id"
	return &id, nil
}
func (m *mockRepository) GetAllMuscularGroups() ([]*dto.MuscularGroupResponseDTO, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.groups, nil
}
func (m *mockRepository) GetMuscularGroupByID(id string) (*dto.MuscularGroupResponseDTO, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	for _, g := range m.groups {
		if g.ID == id {
			return g, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockRepository) UpdateMuscularGroup(id string, mg *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return &dto.MuscularGroupResponseDTO{ID: id, Name: "Updated"}, nil
}

func (m *mockRepository) DeleteMuscularGroup(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return nil
}

func (m *mockRepository) GetMuscularGroupByName(name string) (*dto.MuscularGroupResponseDTO, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}
	for _, g := range m.groups {
		if g.Name == name {
			return g, nil
		}
	}
	return nil, errors.New("not found")
}

func TestCreateMuscularGroup(t *testing.T) {
	mock := &mockRepository{groups: []*dto.MuscularGroupResponseDTO{{ID: "1", Name: "Chest"}}}
	svc := service.NewMuscularGroupService(mock)
	mg := &dto.CreateMuscularGroupDTO{Name: "Back"}
	id, err := svc.CreateMuscularGroup(mg)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "new-id", *id)

	// Duplicate name
	mg.Name = "Chest"
	_, err = svc.CreateMuscularGroup(mg)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "CONFLICT", apiErr.Code)

	// Repo error
	mock.getAllErr = errors.New("fail")
	_, err = svc.CreateMuscularGroup(&dto.CreateMuscularGroupDTO{Name: "Any"})
	assert.Error(t, err)
	mock.getAllErr = nil
	mock.createErr = errors.New("fail")
	_, err = svc.CreateMuscularGroup(&dto.CreateMuscularGroupDTO{Name: "Unique"})
	assert.Error(t, err)
}

func TestGetMuscularGroupByID(t *testing.T) {
	mock := &mockRepository{groups: []*dto.MuscularGroupResponseDTO{{ID: "1", Name: "Chest"}}}
	svc := service.NewMuscularGroupService(mock)
	mg, err := svc.GetMuscularGroupByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "Chest", mg.Name)

	mock.getByIDErr = errors.New("fail")
	_, err = svc.GetMuscularGroupByID("2")
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "NOT_FOUND", apiErr.Code)
}

func TestGetAllMuscularGroups(t *testing.T) {
	mock := &mockRepository{groups: []*dto.MuscularGroupResponseDTO{{ID: "1", Name: "Chest"}}}
	svc := service.NewMuscularGroupService(mock)
	groups, err := svc.GetAllMuscularGroups()
	assert.NoError(t, err)
	assert.Len(t, groups, 1)

	mock.getAllErr = errors.New("fail")
	_, err = svc.GetAllMuscularGroups()
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "INTERNAL_ERROR", apiErr.Code)
}

func TestUpdateMuscularGroup(t *testing.T) {
	mock := &mockRepository{groups: []*dto.MuscularGroupResponseDTO{{ID: "1", Name: "OldName"}, {ID: "2", Name: "Other"}}}
	svc := service.NewMuscularGroupService(mock)
	newName := "NewName"
	update := &dto.UpdateMuscularGroupDTO{Name: &newName}

	// Success
	mg, err := svc.UpdateMuscularGroup("1", update)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", mg.Name)

	// Not found
	mock.getByIDErr = errors.New("not found")
	_, err = svc.UpdateMuscularGroup("3", update)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "NOT_FOUND", apiErr.Code)
	mock.getByIDErr = nil

	// Duplicate name
	update.Name = strPtr("Other")
	_, err = svc.UpdateMuscularGroup("1", update)
	assert.Error(t, err)
	apiErr, ok = err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "CONFLICT", apiErr.Code)

	// Update error
	update.Name = strPtr("NewName")
	mock.updateErr = errors.New("fail")
	_, err = svc.UpdateMuscularGroup("1", update)
	assert.Error(t, err)
	apiErr, ok = err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "INTERNAL_ERROR", apiErr.Code)
	mock.updateErr = nil
}

func TestDeleteMuscularGroup(t *testing.T) {
	mock := &mockRepository{groups: []*dto.MuscularGroupResponseDTO{{ID: "1", Name: "Chest"}}}
	svc := service.NewMuscularGroupService(mock)
	err := svc.DeleteMuscularGroup("1")
	assert.NoError(t, err)

	mock.getByIDErr = errors.New("fail")
	err = svc.DeleteMuscularGroup("2")
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "NOT_FOUND", apiErr.Code)
	mock.getByIDErr = nil

	mock.deleteErr = errors.New("fail")
	err = svc.DeleteMuscularGroup("1")
	assert.Error(t, err)
	apiErr, ok = err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "INTERNAL_ERROR", apiErr.Code)
}

func strPtr(s string) *string { return &s }
