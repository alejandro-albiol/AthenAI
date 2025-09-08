package service

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	CreateFn    func(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error)
	GetByIDFn   func(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error)
	GetByNameFn func(gymID string, name string) (*dto.ResponseCustomEquipmentDTO, error)
	ListFn      func(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error)
	UpdateFn    func(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error
	DeleteFn    func(gymID, id string) error
}

// GetByName implements interfaces.CustomEquipmentRepository.
func (m *mockRepo) GetByName(gymID string, name string) (*dto.ResponseCustomEquipmentDTO, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(gymID, name)
	}
	return nil, nil
}

func (m *mockRepo) Create(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error) {
	return m.CreateFn(gymID, equipment)
}
func (m *mockRepo) GetByID(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error) {
	return m.GetByIDFn(gymID, id)
}
func (m *mockRepo) List(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
	return m.ListFn(gymID)
}
func (m *mockRepo) Update(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error {
	return m.UpdateFn(gymID, equipment)
}
func (m *mockRepo) Delete(gymID, id string) error {
	return m.DeleteFn(gymID, id)
}

func TestCreateCustomEquipmentService(t *testing.T) {
	gymID := "tenant_schema"
	repo := &mockRepo{
		CreateFn: func(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error) {
			if equipment.Name == "" {
				return nil, errors.New("name required")
			}
			id := "eq-1"
			return &id, nil
		},
		GetByNameFn: func(gymID string, name string) (*dto.ResponseCustomEquipmentDTO, error) {
			return nil, nil
		},
	}
	svc := NewCustomEquipmentService(repo)
	equipment := &dto.CreateCustomEquipmentDTO{
		CreatedBy:   "user123",
		Name:        "Dumbbell",
		Description: "A dumbbell",
		Category:    "weight",
		IsActive:    true,
	}
	id, err := svc.CreateCustomEquipment(gymID, equipment)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	badEquipment := &dto.CreateCustomEquipmentDTO{CreatedBy: "user123"}
	id, err = svc.CreateCustomEquipment(gymID, badEquipment)
	assert.Error(t, err)
	assert.Nil(t, id)
}

func TestCreateCustomEquipmentService_Duplicate(t *testing.T) {
	gymID := "tenant_schema"
	repo := &mockRepo{
		CreateFn: func(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error) {
			id := "eq-1"
			return &id, nil
		},
		GetByNameFn: func(gymID string, name string) (*dto.ResponseCustomEquipmentDTO, error) {
			return &dto.ResponseCustomEquipmentDTO{ID: "eq-1", Name: name}, nil
		},
	}
	svc := NewCustomEquipmentService(repo)
	equipment := &dto.CreateCustomEquipmentDTO{
		CreatedBy:   "user123",
		Name:        "Dumbbell",
		Description: "A dumbbell",
		Category:    "weight",
		IsActive:    true,
	}
	id, err := svc.CreateCustomEquipment(gymID, equipment)
	assert.Error(t, err)
	assert.Nil(t, id)
}

func TestGetCustomEquipmentByIDService(t *testing.T) {
	gymID := "tenant_schema"
	repo := &mockRepo{
		GetByIDFn: func(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error) {
			if id == "notfound" {
				return nil, errors.New("not found")
			}
			return &dto.ResponseCustomEquipmentDTO{ID: id, Name: "Dumbbell"}, nil
		},
	}
	svc := NewCustomEquipmentService(repo)
	result, err := svc.GetCustomEquipmentByID(gymID, "eq-1")
	assert.NoError(t, err)
	assert.Equal(t, "eq-1", result.ID)

	_, err = svc.GetCustomEquipmentByID(gymID, "notfound")
	assert.Error(t, err)
}

func TestListCustomEquipmentService(t *testing.T) {
	gymID := "tenant_schema"
	repo := &mockRepo{
		ListFn: func(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
			return []*dto.ResponseCustomEquipmentDTO{{ID: "eq-1", Name: "Dumbbell"}}, nil
		},
	}
	svc := NewCustomEquipmentService(repo)
	result, err := svc.ListCustomEquipment(gymID)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "eq-1", result[0].ID)
}

func TestUpdateCustomEquipmentService(t *testing.T) {
	gymID := "tenant_schema"
	repo := &mockRepo{
		UpdateFn: func(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error {
			if equipment.ID == "notfound" {
				return errors.New("not found")
			}
			return nil
		},
	}
	svc := NewCustomEquipmentService(repo)
	update := &dto.UpdateCustomEquipmentDTO{
		ID:          "eq-1",
		Name:        ptr("Barbell"),
		Description: ptr("A barbell"),
		Category:    ptr("weight"),
		IsActive:    ptr(true),
	}
	err := svc.UpdateCustomEquipment(gymID, update)
	assert.NoError(t, err)

	badUpdate := &dto.UpdateCustomEquipmentDTO{ID: "notfound"}
	err = svc.UpdateCustomEquipment(gymID, badUpdate)
	assert.Error(t, err)
}

func TestDeleteCustomEquipmentService(t *testing.T) {
	gymID := "tenant_schema"
	repo := &mockRepo{
		DeleteFn: func(gymID, id string) error {
			if id == "notfound" {
				return errors.New("not found")
			}
			return nil
		},
	}
	svc := NewCustomEquipmentService(repo)
	err := svc.DeleteCustomEquipment(gymID, "eq-1")
	assert.NoError(t, err)

	err = svc.DeleteCustomEquipment(gymID, "notfound")
	assert.Error(t, err)
}

func ptr[T any](v T) *T { return &v }
