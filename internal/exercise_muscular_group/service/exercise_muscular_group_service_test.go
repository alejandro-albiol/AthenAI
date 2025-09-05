package service

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
)

type mockRepository struct {
	CreateLinkFunc                func(link *dto.ExerciseMuscularGroup) (*string, error)
	DeleteLinkFunc                func(id string) error
	RemoveAllLinksForExerciseFunc func(exerciseID string) error
	FindByIDFunc                  func(id string) (*dto.ExerciseMuscularGroup, error)
	FindByExerciseIDFunc          func(exerciseID string) ([]*dto.ExerciseMuscularGroup, error)
	FindByMuscularGroupIDFunc     func(muscularGroupID string) ([]*dto.ExerciseMuscularGroup, error)
}

func (m *mockRepository) CreateLink(link *dto.ExerciseMuscularGroup) (*string, error) {
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
func (m *mockRepository) FindByID(id string) (*dto.ExerciseMuscularGroup, error) {
	return m.FindByIDFunc(id)
}
func (m *mockRepository) FindByExerciseID(exerciseID string) ([]*dto.ExerciseMuscularGroup, error) {
	if m.FindByExerciseIDFunc != nil {
		return m.FindByExerciseIDFunc(exerciseID)
	}
	return nil, nil
}
func (m *mockRepository) FindByMuscularGroupID(muscularGroupID string) ([]*dto.ExerciseMuscularGroup, error) {
	if m.FindByMuscularGroupIDFunc != nil {
		return m.FindByMuscularGroupIDFunc(muscularGroupID)
	}
	return nil, nil
}

func TestExerciseMuscularGroupService_CreateLink(t *testing.T) {
	mockRepo := &mockRepository{
		CreateLinkFunc: func(link *dto.ExerciseMuscularGroup) (*string, error) {
			id := "ex1"
			return &id, nil
		},
	}
	service := NewExerciseMuscularGroupService(mockRepo)

	link := &dto.ExerciseMuscularGroup{ExerciseID: "ex1", MuscularGroupID: "mg1"}
	res, err := service.CreateLink(link)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res == nil || *res != "ex1" {
		t.Errorf("expected id 'ex1', got %v", res)
	}

	// Error case
	mockRepo.CreateLinkFunc = func(link *dto.ExerciseMuscularGroup) (*string, error) {
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

func TestExerciseMuscularGroupService_DeleteLink(t *testing.T) {
	mockRepo := &mockRepository{
		DeleteLinkFunc: func(id string) error {
			return nil
		},
	}
	service := NewExerciseMuscularGroupService(mockRepo)

	err := service.DeleteLink("ex1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Error case
	mockRepo.DeleteLinkFunc = func(id string) error {
		return errors.New("repo error")
	}
	err = service.DeleteLink("ex1")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestExerciseMuscularGroupService_GetLinkByID(t *testing.T) {
	mockRepo := &mockRepository{
		FindByIDFunc: func(id string) (*dto.ExerciseMuscularGroup, error) {
			return &dto.ExerciseMuscularGroup{ExerciseID: "ex1", MuscularGroupID: "mg1"}, nil
		},
	}
	service := NewExerciseMuscularGroupService(mockRepo)

	link, err := service.GetLinkByID("ex1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if link == nil || link.ExerciseID != "ex1" || link.MuscularGroupID != "mg1" {
		t.Errorf("unexpected link result: %+v", link)
	}

	// Error case
	mockRepo.FindByIDFunc = func(id string) (*dto.ExerciseMuscularGroup, error) {
		return nil, errors.New("repo error")
	}
	link, err = service.GetLinkByID("ex1")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if link != nil {
		t.Errorf("expected nil link, got %+v", link)
	}
}

func TestExerciseMuscularGroupService_RemoveAllLinksForExercise(t *testing.T) {
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
	service := NewExerciseMuscularGroupService(mockRepo)

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

func TestExerciseMuscularGroupService_GetLinksByExerciseID(t *testing.T) {
	mockRepo := &mockRepository{
		FindByExerciseIDFunc: func(exerciseID string) ([]*dto.ExerciseMuscularGroup, error) {
			if exerciseID == "fail" {
				return nil, errors.New("repo error")
			}
			return []*dto.ExerciseMuscularGroup{{ExerciseID: exerciseID, MuscularGroupID: "mg1"}}, nil
		},
	}
	service := NewExerciseMuscularGroupService(mockRepo)

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

func TestExerciseMuscularGroupService_GetLinksByMuscularGroupID(t *testing.T) {
	mockRepo := &mockRepository{
		FindByMuscularGroupIDFunc: func(muscularGroupID string) ([]*dto.ExerciseMuscularGroup, error) {
			if muscularGroupID == "fail" {
				return nil, errors.New("repo error")
			}
			return []*dto.ExerciseMuscularGroup{{ExerciseID: "ex1", MuscularGroupID: muscularGroupID}}, nil
		},
	}
	service := NewExerciseMuscularGroupService(mockRepo)

	// Success case
	links, err := service.GetLinksByMuscularGroupID("mg1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(links) != 1 || links[0].MuscularGroupID != "mg1" {
		t.Errorf("unexpected links: %+v", links)
	}

	// Error case
	links, err = service.GetLinksByMuscularGroupID("fail")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if links != nil {
		t.Errorf("expected nil links, got %+v", links)
	}
}
