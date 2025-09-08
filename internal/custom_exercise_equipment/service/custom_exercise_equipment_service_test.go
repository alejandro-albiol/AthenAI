package service

import (
	"errors"
	"testing"

	customEquipmentDTO "github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"
	equipmentDTO "github.com/alejandro-albiol/athenai/internal/equipment/dto"
	"github.com/stretchr/testify/assert"
)

type mockCustomEquipmentRepository struct{}

func (f *mockCustomEquipmentRepository) GetByID(gymID, equipmentID string) (*customEquipmentDTO.ResponseCustomEquipmentDTO, error) {
	return &customEquipmentDTO.ResponseCustomEquipmentDTO{ID: equipmentID}, nil
}

// Add the missing Create method to satisfy the interface
func (f *mockCustomEquipmentRepository) Create(gymID string, equipment *customEquipmentDTO.CreateCustomEquipmentDTO) (*string, error) {
	id := "fake-id"
	return &id, nil
}

// Add the missing Delete method to satisfy the interface
func (f *mockCustomEquipmentRepository) Delete(gymID, equipmentID string) error {
	return nil
}

func (f *mockCustomEquipmentRepository) List(gymID string) ([]*customEquipmentDTO.ResponseCustomEquipmentDTO, error) {
	return []*customEquipmentDTO.ResponseCustomEquipmentDTO{}, nil
}
func (f *mockCustomEquipmentRepository) Update(gymID string, equipment *customEquipmentDTO.UpdateCustomEquipmentDTO) error {
	return nil
}
func (f *mockCustomEquipmentRepository) GetByName(gymID, name string) (*customEquipmentDTO.ResponseCustomEquipmentDTO, error) {
	return &customEquipmentDTO.ResponseCustomEquipmentDTO{ID: "mock-id", Name: name}, nil
}

// Satisfy EquipmentRepository interface
func (f *mockPublicEquipmentRepository) GetAllEquipment() ([]*equipmentDTO.EquipmentResponseDTO, error) {
	return []*equipmentDTO.EquipmentResponseDTO{}, nil
}

type mockPublicEquipmentRepository struct{}

func (f *mockPublicEquipmentRepository) GetEquipmentByID(equipmentID string) (*equipmentDTO.EquipmentResponseDTO, error) {
	return &equipmentDTO.EquipmentResponseDTO{ID: equipmentID}, nil
}

func (f *mockPublicEquipmentRepository) CreateEquipment(equipment *equipmentDTO.EquipmentCreationDTO) (*string, error) {
	id := "fake-equip-id"
	return &id, nil
}
func (f *mockPublicEquipmentRepository) UpdateEquipment(id string, update *equipmentDTO.EquipmentUpdateDTO) (*equipmentDTO.EquipmentResponseDTO, error) {
	return &equipmentDTO.EquipmentResponseDTO{ID: id}, nil
}
func (f *mockPublicEquipmentRepository) DeleteEquipment(id string) error {
	return nil
}
func (f *mockPublicEquipmentRepository) ListEquipment() ([]interface{}, error) {
	return nil, nil
}

type mockRepo struct {
	CreateLinkFn                func(gymID string, link *dto.CustomExerciseEquipment) (*string, error)
	DeleteLinkFn                func(gymID, id string) error
	FindByIDFn                  func(gymID, id string) (*dto.CustomExerciseEquipment, error)
	FindByCustomExerciseIDFn    func(gymID, customExerciseID string) ([]*dto.CustomExerciseEquipment, error)
	FindByEquipmentIDFn         func(gymID, equipmentID string) ([]*dto.CustomExerciseEquipment, error)
	RemoveAllLinksForExerciseFn func(gymID, customExerciseID string) error
}

func (m *mockRepo) CreateLink(gymID string, link *dto.CustomExerciseEquipment) (*string, error) {
	return m.CreateLinkFn(gymID, link)
}
func (m *mockRepo) DeleteLink(gymID, id string) error {
	return m.DeleteLinkFn(gymID, id)
}

func (m *mockRepo) FindByID(gymID, id string) (*dto.CustomExerciseEquipment, error) {
	return m.FindByIDFn(gymID, id)
}
func (m *mockRepo) FindByCustomExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseEquipment, error) {
	return m.FindByCustomExerciseIDFn(gymID, customExerciseID)
}
func (m *mockRepo) FindByEquipmentID(gymID, equipmentID string) ([]*dto.CustomExerciseEquipment, error) {
	if m.FindByEquipmentIDFn != nil {
		return m.FindByEquipmentIDFn(gymID, equipmentID)
	}
	return nil, nil
}

func (m *mockRepo) RemoveAllLinksForExercise(gymID, customExerciseID string) error {
	if m.RemoveAllLinksForExerciseFn != nil {
		return m.RemoveAllLinksForExerciseFn(gymID, customExerciseID)
	}
	return nil
}

func TestCreateLinkService(t *testing.T) {
	repo := &mockRepo{
		CreateLinkFn: func(gymID string, link *dto.CustomExerciseEquipment) (*string, error) {
			if link.CustomExerciseID == "" || link.EquipmentID == "" {
				return nil, errors.New("missing fields")
			}
			id := "link-1"
			return &id, nil
		},
	}
	// Mock customEquipmentRepo and publicEquipmentRepo to always allow equipment
	svc := &CustomExerciseEquipmentService{
		repository:          repo,
		customEquipmentRepo: &mockCustomEquipmentRepository{},
		publicEquipmentRepo: &mockPublicEquipmentRepository{},
	}
	link := &dto.CustomExerciseEquipment{CustomExerciseID: "ex-1", EquipmentID: "eq-1"}
	id, err := svc.CreateLink("tenant1", link)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "link-1", *id)
}

func TestFindByIDService(t *testing.T) {
	repo := &mockRepo{
		FindByIDFn: func(gymID, id string) (*dto.CustomExerciseEquipment, error) {
			if id == "notfound" {
				return nil, errors.New("not found")
			}
			return &dto.CustomExerciseEquipment{ID: id, CustomExerciseID: "ex-1", EquipmentID: "eq-1"}, nil
		},
	}
	svc := &CustomExerciseEquipmentService{repository: repo}
	link, err := svc.FindByID("tenant1", "link-1")
	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, "link-1", link.ID)

	_, err = svc.FindByID("tenant1", "notfound")
	assert.Error(t, err)
}
