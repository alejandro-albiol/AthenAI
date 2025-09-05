package service

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
)

type mockRepository struct {
CreateLinkFunc                func(link *dto.ExerciseEquipment) (*string, error)
DeleteLinkFunc                func(id string) error
RemoveAllLinksForExerciseFunc func(exerciseID string) error
FindByIDFunc                  func(id string) (*dto.ExerciseEquipment, error)
FindByExerciseIDFunc          func(exerciseID string) ([]*dto.ExerciseEquipment, error)
FindByEquipmentIDFunc         func(equipmentID string) ([]*dto.ExerciseEquipment, error)
}

func (m *mockRepository) CreateLink(link *dto.ExerciseEquipment) (*string, error) {
	return m.CreateLinkFunc(link)
}
func (m *mockRepository) DeleteLink(id string) error {
	return m.DeleteLinkFunc(id)
}
func (m *mockRepository) RemoveAllLinksForExercise(exerciseID string) error {
       if m.RemoveAllLinksForExerciseFunc != nil {
	       return m.RemoveAllLinksForExerciseFunc(exerciseID)
       }
       return nil
}
func (m *mockRepository) FindByID(id string) (*dto.ExerciseEquipment, error) {
	return m.FindByIDFunc(id)
}
func (m *mockRepository) FindByExerciseID(exerciseID string) ([]*dto.ExerciseEquipment, error) {
       if m.FindByExerciseIDFunc != nil {
	       return m.FindByExerciseIDFunc(exerciseID)
       }
       return nil, nil
}
func (m *mockRepository) FindByEquipmentID(equipmentID string) ([]*dto.ExerciseEquipment, error) {
       if m.FindByEquipmentIDFunc != nil {
	       return m.FindByEquipmentIDFunc(equipmentID)
       }
       return nil, nil
}

func TestExerciseEquipmentService_CreateLink(t *testing.T) {
	mockRepo := &mockRepository{
		CreateLinkFunc: func(link *dto.ExerciseEquipment) (*string, error) {
			id := "id1"
			return &id, nil
		},
	}
	service := NewExerciseEquipmentService(mockRepo)

	link := &dto.ExerciseEquipment{ExerciseID: "ex1", EquipmentID: "eq1"}
	res, err := service.CreateLink(link)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res == nil || *res != "id1" {
		t.Errorf("expected id 'id1', got %v", res)
	}

	// Error case
	mockRepo.CreateLinkFunc = func(link *dto.ExerciseEquipment) (*string, error) {
		return nil, errors.New("repo error")
	}
	res, err = service.CreateLink(link)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if res != nil {
		t.Errorf("expected nil id, got %v", res)
	}
}

func TestExerciseEquipmentService_DeleteLink(t *testing.T) {
	mockRepo := &mockRepository{
		DeleteLinkFunc: func(id string) error {
			return nil
		},
	}
	service := NewExerciseEquipmentService(mockRepo)

	err := service.DeleteLink("id1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Error case
	mockRepo.DeleteLinkFunc = func(id string) error {
		return errors.New("repo error")
	}
	err = service.DeleteLink("id1")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestExerciseEquipmentService_GetLinkByID(t *testing.T) {
	mockRepo := &mockRepository{
		FindByIDFunc: func(id string) (*dto.ExerciseEquipment, error) {
			return &dto.ExerciseEquipment{ExerciseID: "ex1", EquipmentID: "eq1"}, nil
		},
	}
	service := NewExerciseEquipmentService(mockRepo)

	link, err := service.GetLinkByID("id1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if link == nil || link.ExerciseID != "ex1" || link.EquipmentID != "eq1" {
		t.Errorf("unexpected link result: %+v", link)
	}

	// Error case
	mockRepo.FindByIDFunc = func(id string) (*dto.ExerciseEquipment, error) {
		return nil, errors.New("repo error")
	}
	link, err = service.GetLinkByID("id1")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if link != nil {
		t.Errorf("expected nil link, got %+v", link)
	}
}


func TestExerciseEquipmentService_RemoveAllLinksForExercise(t *testing.T) {
       called := false
       mockRepo := &mockRepository{
	       RemoveAllLinksForExerciseFunc: func(exerciseID string) error {
		       called = true
		       if exerciseID == "fail" {
			       return errors.New("repo error")
		       }
		       return nil
	       },
       }
       service := NewExerciseEquipmentService(mockRepo)

       // Success case
       err := service.RemoveAllLinksForExercise("ex1")
       if err != nil {
	       t.Errorf("expected no error, got %v", err)
       }
       if !called {
	       t.Error("expected RemoveAllLinksForExercise to be called")
       }

       // Error case
       called = false
       err = service.RemoveAllLinksForExercise("fail")
       if err == nil {
	       t.Error("expected error, got nil")
       }
       if !called {
	       t.Error("expected RemoveAllLinksForExercise to be called")
       }
}

func TestExerciseEquipmentService_GetLinksByExerciseID(t *testing.T) {
       mockRepo := &mockRepository{
	       FindByExerciseIDFunc: func(exerciseID string) ([]*dto.ExerciseEquipment, error) {
		       if exerciseID == "fail" {
			       return nil, errors.New("repo error")
		       }
		       return []*dto.ExerciseEquipment{{ExerciseID: exerciseID, EquipmentID: "eq1"}}, nil
	       },
       }
       service := NewExerciseEquipmentService(mockRepo)

       // Success case
       links, err := service.GetLinksByExerciseID("ex1")
       if err != nil {
	       t.Errorf("expected no error, got %v", err)
       }
       if len(links) != 1 || links[0].ExerciseID != "ex1" {
	       t.Errorf("unexpected links: %+v", links)
       }

       // Error case
       links, err = service.GetLinksByExerciseID("fail")
       if err == nil {
	       t.Error("expected error, got nil")
       }
       if links != nil {
	       t.Errorf("expected nil links, got %+v", links)
       }
}

func TestExerciseEquipmentService_GetLinksByEquipmentID(t *testing.T) {
       mockRepo := &mockRepository{
	       FindByEquipmentIDFunc: func(equipmentID string) ([]*dto.ExerciseEquipment, error) {
		       if equipmentID == "fail" {
			       return nil, errors.New("repo error")
		       }
		       if equipmentID == "empty" {
			       return []*dto.ExerciseEquipment{}, nil
		       }
		       return []*dto.ExerciseEquipment{{ExerciseID: "ex1", EquipmentID: equipmentID}}, nil
	       },
       }
       service := NewExerciseEquipmentService(mockRepo)

       // Success case
       links, err := service.GetLinksByEquipmentID("eq1")
       if err != nil {
	       t.Errorf("expected no error, got %v", err)
       }
       if len(links) != 1 || links[0].EquipmentID != "eq1" {
	       t.Errorf("unexpected links: %+v", links)
       }

       // Not found case
       links, err = service.GetLinksByEquipmentID("empty")
       if err == nil {
	       t.Error("expected error, got nil")
       }
       if links != nil {
	       t.Errorf("expected nil links, got %+v", links)
       }

       // Error case
       links, err = service.GetLinksByEquipmentID("fail")
       if err == nil {
	       t.Error("expected error, got nil")
       }
       if links != nil {
	       t.Errorf("expected nil links, got %+v", links)
       }
}
