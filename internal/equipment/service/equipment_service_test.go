package service

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/equipment/dto"
)

type mockEquipmentRepository struct {
	CreateEquipmentFunc  func(equipment *dto.EquipmentCreationDTO) (*string, error)
	GetEquipmentByIDFunc func(id string) (*dto.EquipmentResponseDTO, error)
	GetAllEquipmentFunc  func() ([]*dto.EquipmentResponseDTO, error)
	UpdateEquipmentFunc  func(id string, update *dto.EquipmentUpdateDTO) (*dto.EquipmentResponseDTO, error)
	DeleteEquipmentFunc  func(id string) error
}

func (m *mockEquipmentRepository) CreateEquipment(equipment *dto.EquipmentCreationDTO) (*string, error) {
	return m.CreateEquipmentFunc(equipment)
}
func (m *mockEquipmentRepository) GetEquipmentByID(id string) (*dto.EquipmentResponseDTO, error) {
	return m.GetEquipmentByIDFunc(id)
}
func (m *mockEquipmentRepository) GetAllEquipment() ([]*dto.EquipmentResponseDTO, error) {
	return m.GetAllEquipmentFunc()
}
func (m *mockEquipmentRepository) UpdateEquipment(id string, update *dto.EquipmentUpdateDTO) (*dto.EquipmentResponseDTO, error) {
	return m.UpdateEquipmentFunc(id, update)
}
func (m *mockEquipmentRepository) DeleteEquipment(id string) error {
	return m.DeleteEquipmentFunc(id)
}

func TestEquipmentService_CreateEquipment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetAllEquipmentFunc: func() ([]*dto.EquipmentResponseDTO, error) {
				return []*dto.EquipmentResponseDTO{}, nil
			},
			CreateEquipmentFunc: func(e *dto.EquipmentCreationDTO) (*string, error) {
				id := "id1"
				return &id, nil
			},
		}
		service := NewEquipmentService(mockRepo)
		input := &dto.EquipmentCreationDTO{Name: "Dumbbell", Description: "A dumbbell", Category: "free_weights"}
		id, err := service.CreateEquipment(input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if id == nil || *id != "id1" {
			t.Errorf("expected id 'id1', got %v", id)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetAllEquipmentFunc: func() ([]*dto.EquipmentResponseDTO, error) {
				return []*dto.EquipmentResponseDTO{{Name: "Dumbbell"}}, nil
			},
		}
		service := NewEquipmentService(mockRepo)
		input := &dto.EquipmentCreationDTO{Name: "Dumbbell"}
		id, err := service.CreateEquipment(input)
		if err == nil || id != nil {
			t.Error("expected conflict error for duplicate name")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetAllEquipmentFunc: func() ([]*dto.EquipmentResponseDTO, error) {
				return nil, errors.New("db error")
			},
		}
		service := NewEquipmentService(mockRepo)
		input := &dto.EquipmentCreationDTO{Name: "Barbell"}
		id, err := service.CreateEquipment(input)
		if err == nil || id != nil {
			t.Error("expected error for repo failure")
		}
	})
}

func TestEquipmentService_GetEquipmentByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return &dto.EquipmentResponseDTO{ID: id, Name: "Dumbbell"}, nil
			},
		}
		service := NewEquipmentService(mockRepo)
		res, err := service.GetEquipmentByID("id1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res == nil || res.ID != "id1" {
			t.Errorf("expected id 'id1', got %v", res)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return nil, errors.New("not found")
			},
		}
		service := NewEquipmentService(mockRepo)
		res, err := service.GetEquipmentByID("id404")
		if err == nil || res != nil {
			t.Error("expected not found error")
		}
	})
}

func TestEquipmentService_GetAllEquipment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetAllEquipmentFunc: func() ([]*dto.EquipmentResponseDTO, error) {
				return []*dto.EquipmentResponseDTO{{ID: "id1"}}, nil
			},
		}
		service := NewEquipmentService(mockRepo)
		res, err := service.GetAllEquipment()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(res) != 1 || res[0].ID != "id1" {
			t.Errorf("expected one result with id 'id1', got %v", res)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetAllEquipmentFunc: func() ([]*dto.EquipmentResponseDTO, error) {
				return nil, errors.New("db error")
			},
		}
		service := NewEquipmentService(mockRepo)
		res, err := service.GetAllEquipment()
		if err == nil || res != nil {
			t.Error("expected error for repo failure")
		}
	})
}

func TestEquipmentService_UpdateEquipment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return &dto.EquipmentResponseDTO{ID: id}, nil
			},
			UpdateEquipmentFunc: func(id string, update *dto.EquipmentUpdateDTO) (*dto.EquipmentResponseDTO, error) {
				return &dto.EquipmentResponseDTO{ID: id, Name: "Updated"}, nil
			},
		}
		service := NewEquipmentService(mockRepo)
		upd := &dto.EquipmentUpdateDTO{}
		res, err := service.UpdateEquipment("id1", upd)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res == nil || res.ID != "id1" {
			t.Errorf("expected id 'id1', got %v", res)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return nil, errors.New("not found")
			},
		}
		service := NewEquipmentService(mockRepo)
		upd := &dto.EquipmentUpdateDTO{}
		res, err := service.UpdateEquipment("id404", upd)
		if err == nil || res != nil {
			t.Error("expected not found error")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return &dto.EquipmentResponseDTO{ID: id}, nil
			},
			UpdateEquipmentFunc: func(id string, update *dto.EquipmentUpdateDTO) (*dto.EquipmentResponseDTO, error) {
				return nil, errors.New("db error")
			},
		}
		service := NewEquipmentService(mockRepo)
		upd := &dto.EquipmentUpdateDTO{}
		res, err := service.UpdateEquipment("id1", upd)
		if err == nil || res != nil {
			t.Error("expected error for repo failure")
		}
	})
}

func TestEquipmentService_DeleteEquipment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return &dto.EquipmentResponseDTO{ID: id}, nil
			},
			DeleteEquipmentFunc: func(id string) error {
				return nil
			},
		}
		service := NewEquipmentService(mockRepo)
		err := service.DeleteEquipment("id1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return nil, errors.New("not found")
			},
		}
		service := NewEquipmentService(mockRepo)
		err := service.DeleteEquipment("id404")
		if err == nil {
			t.Error("expected not found error")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := &mockEquipmentRepository{
			GetEquipmentByIDFunc: func(id string) (*dto.EquipmentResponseDTO, error) {
				return &dto.EquipmentResponseDTO{ID: id}, nil
			},
			DeleteEquipmentFunc: func(id string) error {
				return errors.New("db error")
			},
		}
		service := NewEquipmentService(mockRepo)
		err := service.DeleteEquipment("id1")
		if err == nil {
			t.Error("expected error for repo failure")
		}
	})
}
