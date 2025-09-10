package service

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/exercise/dto"

	exerciseEquipmentDTO "github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"

	exerciseMuscularGroupDTO "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
)

type mockRepository struct {
	CreateExerciseFunc              func(ex *dto.ExerciseCreationDTO) (*string, error)
	GetExerciseByNameFunc           func(name string) (*dto.ExerciseResponseDTO, error)
	DeleteExerciseFunc              func(id string) error
	GetExerciseByIDFunc             func(id string) (*dto.ExerciseResponseDTO, error)
	GetExercisesByMuscularGroupFunc func(muscularGroups []string) ([]*dto.ExerciseResponseDTO, error)
	GetExercisesByEquipmentFunc     func(equipment []string) ([]*dto.ExerciseResponseDTO, error)
	GetAllExercisesFunc             func() ([]*dto.ExerciseResponseDTO, error)
	UpdateExerciseFunc              func(id string, exercise *dto.ExerciseUpdateDTO) (*dto.ExerciseResponseDTO, error)
}

func (m *mockRepository) CreateExercise(ex *dto.ExerciseCreationDTO) (*string, error) {
	return m.CreateExerciseFunc(ex)
}
func (m *mockRepository) GetExerciseByName(name string) (*dto.ExerciseResponseDTO, error) {
	return m.GetExerciseByNameFunc(name)
}
func (m *mockRepository) DeleteExercise(id string) error {
	if m.DeleteExerciseFunc != nil {
		return m.DeleteExerciseFunc(id)
	}
	return nil
}
func (m *mockRepository) GetExerciseByID(id string) (*dto.ExerciseResponseDTO, error) {
	if m.GetExerciseByIDFunc != nil {
		return m.GetExerciseByIDFunc(id)
	}
	return nil, nil
}
func (m *mockRepository) GetExercisesByMuscularGroup(muscularGroups []string) ([]*dto.ExerciseResponseDTO, error) {
	if m.GetExercisesByMuscularGroupFunc != nil {
		return m.GetExercisesByMuscularGroupFunc(muscularGroups)
	}
	return nil, nil
}
func (m *mockRepository) GetExercisesByEquipment(equipment []string) ([]*dto.ExerciseResponseDTO, error) {
	if m.GetExercisesByEquipmentFunc != nil {
		return m.GetExercisesByEquipmentFunc(equipment)
	}
	return nil, nil
}
func (m *mockRepository) GetAllExercises() ([]*dto.ExerciseResponseDTO, error) {
	if m.GetAllExercisesFunc != nil {
		return m.GetAllExercisesFunc()
	}
	return nil, nil
}
func (m *mockRepository) UpdateExercise(id string, exercise *dto.ExerciseUpdateDTO) (*dto.ExerciseResponseDTO, error) {
	if m.UpdateExerciseFunc != nil {
		return m.UpdateExerciseFunc(id, exercise)
	}
	return nil, nil
}

// ...implement other methods as needed for further tests

type mockEquipmentService struct {
	CreateLinkFunc                func(link *exerciseEquipmentDTO.ExerciseEquipment) (*string, error)
	DeleteLinkFunc                func(id string) error
	RemoveAllLinksForExerciseFunc func(exerciseID string) error
	GetLinkByIDFunc               func(id string) (*exerciseEquipmentDTO.ExerciseEquipment, error)
	GetLinksByExerciseIDFunc      func(exerciseID string) ([]*exerciseEquipmentDTO.ExerciseEquipment, error)
	GetLinksByEquipmentIDFunc     func(equipmentID string) ([]*exerciseEquipmentDTO.ExerciseEquipment, error)
}

func (m *mockEquipmentService) CreateLink(link *exerciseEquipmentDTO.ExerciseEquipment) (*string, error) {
	return m.CreateLinkFunc(link)
}
func (m *mockEquipmentService) RemoveAllLinksForExercise(exerciseID string) error {
	if m.RemoveAllLinksForExerciseFunc != nil {
		return m.RemoveAllLinksForExerciseFunc(exerciseID)
	}
	return nil
}
func (m *mockEquipmentService) DeleteLink(id string) error {
	if m.DeleteLinkFunc != nil {
		return m.DeleteLinkFunc(id)
	}
	return nil
}
func (m *mockEquipmentService) GetLinkByID(id string) (*exerciseEquipmentDTO.ExerciseEquipment, error) {
	if m.GetLinkByIDFunc != nil {
		return m.GetLinkByIDFunc(id)
	}
	return nil, nil
}
func (m *mockEquipmentService) GetLinksByExerciseID(exerciseID string) ([]*exerciseEquipmentDTO.ExerciseEquipment, error) {
	if m.GetLinksByExerciseIDFunc != nil {
		return m.GetLinksByExerciseIDFunc(exerciseID)
	}
	return nil, nil
}
func (m *mockEquipmentService) GetLinksByEquipmentID(equipmentID string) ([]*exerciseEquipmentDTO.ExerciseEquipment, error) {
	if m.GetLinksByEquipmentIDFunc != nil {
		return m.GetLinksByEquipmentIDFunc(equipmentID)
	}
	return nil, nil
}

type mockMuscularGroupService struct {
	CreateLinkFunc                func(link *exerciseMuscularGroupDTO.ExerciseMuscularGroup) (*string, error)
	DeleteLinkFunc                func(id string) error
	RemoveAllLinksForExerciseFunc func(exerciseID string) error
	GetLinkByIDFunc               func(id string) (*exerciseMuscularGroupDTO.ExerciseMuscularGroup, error)
	GetLinksByExerciseIDFunc      func(exerciseID string) ([]*exerciseMuscularGroupDTO.ExerciseMuscularGroup, error)
	GetLinksByMuscularGroupIDFunc func(muscularGroupID string) ([]*exerciseMuscularGroupDTO.ExerciseMuscularGroup, error)
}

func (m *mockMuscularGroupService) CreateLink(link *exerciseMuscularGroupDTO.ExerciseMuscularGroup) (*string, error) {
	return m.CreateLinkFunc(link)
}
func (m *mockMuscularGroupService) RemoveAllLinksForExercise(exerciseID string) error {
	if m.RemoveAllLinksForExerciseFunc != nil {
		return m.RemoveAllLinksForExerciseFunc(exerciseID)
	}
	return nil
}
func (m *mockMuscularGroupService) DeleteLink(id string) error {
	if m.DeleteLinkFunc != nil {
		return m.DeleteLinkFunc(id)
	}
	return nil
}
func (m *mockMuscularGroupService) GetLinkByID(id string) (*exerciseMuscularGroupDTO.ExerciseMuscularGroup, error) {
	if m.GetLinkByIDFunc != nil {
		return m.GetLinkByIDFunc(id)
	}
	return nil, nil
}
func (m *mockMuscularGroupService) GetLinksByExerciseID(exerciseID string) ([]*exerciseMuscularGroupDTO.ExerciseMuscularGroup, error) {
	if m.GetLinksByExerciseIDFunc != nil {
		return m.GetLinksByExerciseIDFunc(exerciseID)
	}
	return nil, nil
}
func (m *mockMuscularGroupService) GetLinksByMuscularGroupID(muscularGroupID string) ([]*exerciseMuscularGroupDTO.ExerciseMuscularGroup, error) {
	if m.GetLinksByMuscularGroupIDFunc != nil {
		return m.GetLinksByMuscularGroupIDFunc(muscularGroupID)
	}
	return nil, nil
}

func TestExerciseService_CreateExercise(t *testing.T) {
	repo := &mockRepository{
		CreateExerciseFunc: func(ex *dto.ExerciseCreationDTO) (*string, error) {
			id := "id1"
			return &id, nil
		},
		GetExerciseByNameFunc: func(name string) (*dto.ExerciseResponseDTO, error) {
			return &dto.ExerciseResponseDTO{}, errors.New("not found")
		},
	}
	equipmentSvc := &mockEquipmentService{
		CreateLinkFunc: func(link *exerciseEquipmentDTO.ExerciseEquipment) (*string, error) {
			id := "link1"
			return &id, nil
		},
	}
	muscularGroupSvc := &mockMuscularGroupService{
		CreateLinkFunc: func(link *exerciseMuscularGroupDTO.ExerciseMuscularGroup) (*string, error) {
			id := "link2"
			return &id, nil
		},
	}
	service := NewExerciseService(repo, equipmentSvc, muscularGroupSvc)

	ex := &dto.ExerciseCreationDTO{
		Name:            "Pushup",
		DifficultyLevel: "beginner",
		ExerciseType:    "strength",
		Synonyms:        []string{"press-up"},
		Equipment:       []string{"eq1"},
		MuscularGroups:  []string{"mg1"},
		Instructions:    "Do a pushup",
		CreatedBy:       "tester",
	}
	id, err := service.CreateExercise(ex)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if id == nil || *id != "id1" {
		t.Errorf("expected id 'id1', got %v", id)
	}

	// Error: duplicate name
	repo.GetExerciseByNameFunc = func(name string) (*dto.ExerciseResponseDTO, error) {
		return &dto.ExerciseResponseDTO{ID: "exists"}, nil
	}
	_, err = service.CreateExercise(ex)
	if err == nil {
		t.Error("expected error for duplicate name, got nil")
	}

	// Error: equipment link fails
	repo.GetExerciseByNameFunc = func(name string) (*dto.ExerciseResponseDTO, error) {
		return &dto.ExerciseResponseDTO{}, errors.New("not found")
	}
	equipmentSvc.CreateLinkFunc = func(link *exerciseEquipmentDTO.ExerciseEquipment) (*string, error) {
		return nil, errors.New("equipment link error")
	}
	_, err = service.CreateExercise(ex)
	if err == nil {
		t.Error("expected error for equipment link, got nil")
	}

	// Error: muscular group link fails
	equipmentSvc.CreateLinkFunc = func(link *exerciseEquipmentDTO.ExerciseEquipment) (*string, error) {
		id := "link1"
		return &id, nil
	}
	muscularGroupSvc.CreateLinkFunc = func(link *exerciseMuscularGroupDTO.ExerciseMuscularGroup) (*string, error) {
		return nil, errors.New("mg link error")
	}
	_, err = service.CreateExercise(ex)
	if err == nil {
		t.Error("expected error for muscular group link, got nil")
	}
}

func TestExerciseService_DeleteExercise(t *testing.T) {
	repo := &mockRepository{
		DeleteExerciseFunc: func(id string) error {
			if id == "fail" {
				return errors.New("delete error")
			}
			return nil
		},
		GetExerciseByIDFunc: func(id string) (*dto.ExerciseResponseDTO, error) {
			if id == "notfound" {
				return nil, errors.New("not found")
			}
			return &dto.ExerciseResponseDTO{ID: id, Name: "Pushup"}, nil
		},
	}
	equipmentSvc := &mockEquipmentService{
		RemoveAllLinksForExerciseFunc: func(exerciseID string) error {
			return errors.New("equipment remove error")
		},
	}
	muscularGroupSvc := &mockMuscularGroupService{
		RemoveAllLinksForExerciseFunc: func(exerciseID string) error {
			return errors.New("mg remove error")
		},
	}
	service := NewExerciseService(repo, equipmentSvc, muscularGroupSvc)

	// Success
	err := service.DeleteExercise("ok")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Not found
	err = service.DeleteExercise("notfound")
	if err == nil {
		t.Error("expected error for not found, got nil")
	}

	// Repo delete error
	err = service.DeleteExercise("fail")
	if err == nil || err.Error() != "Failed to delete exercise" {
		t.Errorf("expected delete error, got %v", err)
	}
}

func TestExerciseService_GetExerciseByID(t *testing.T) {
	repo := &mockRepository{
		GetExerciseByIDFunc: func(id string) (*dto.ExerciseResponseDTO, error) {
			if id == "notfound" {
				return nil, errors.New("not found")
			}
			return &dto.ExerciseResponseDTO{ID: id, Name: "Pushup"}, nil
		},
	}
	equipmentSvc := &mockEquipmentService{}
	muscularGroupSvc := &mockMuscularGroupService{}
	service := NewExerciseService(repo, equipmentSvc, muscularGroupSvc)

	// Success
	res, err := service.GetExerciseByID("id1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res == nil || res.ID != "id1" {
		t.Errorf("expected id 'id1', got %v", res)
	}

	// Not found
	res, err = service.GetExerciseByID("notfound")
	if err == nil {
		t.Error("expected error for not found, got nil")
	}
	if res != nil {
		t.Errorf("expected nil result for not found, got %v", res)
	}
}

func TestExerciseService_UpdateExercise(t *testing.T) {
	repo := &mockRepository{
		UpdateExerciseFunc: func(id string, ex *dto.ExerciseUpdateDTO) (*dto.ExerciseResponseDTO, error) {
			if id == "fail" {
				return nil, errors.New("update error")
			}
			return &dto.ExerciseResponseDTO{ID: id, Name: "Updated"}, nil
		},
	}
	equipmentSvc := &mockEquipmentService{}
	muscularGroupSvc := &mockMuscularGroupService{}
	service := NewExerciseService(repo, equipmentSvc, muscularGroupSvc)

	// Success
	name := "Updated"
	upd := &dto.ExerciseUpdateDTO{Name: &name}
	res, err := service.UpdateExercise("id1", upd)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res == nil || res.ID != "id1" {
		t.Errorf("expected id 'id1', got %v", res)
	}

	// Update error
	_, err = service.UpdateExercise("fail", upd)
	if err == nil {
		t.Error("expected error for update, got nil")
	}
}

func TestExerciseService_GetAllExercises(t *testing.T) {
	repo := &mockRepository{
		GetAllExercisesFunc: func() ([]*dto.ExerciseResponseDTO, error) {
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	service := NewExerciseService(repo, &mockEquipmentService{}, &mockMuscularGroupService{})
	res, err := service.GetAllExercises()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0].ID != "id1" {
		t.Errorf("expected one result with id 'id1', got %v", res)
	}
}

func TestExerciseService_GetExercisesByMuscularGroup(t *testing.T) {
	repo := &mockRepository{
		GetExercisesByMuscularGroupFunc: func(mg []string) ([]*dto.ExerciseResponseDTO, error) {
			if len(mg) == 1 && mg[0] == "fail" {
				return nil, errors.New("mg error")
			}
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	service := NewExerciseService(repo, &mockEquipmentService{}, &mockMuscularGroupService{})
	res, err := service.GetExercisesByMuscularGroup([]string{"mg1"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0].ID != "id1" {
		t.Errorf("expected one result with id 'id1', got %v", res)
	}
	_, err = service.GetExercisesByMuscularGroup([]string{"fail"})
	if err == nil {
		t.Error("expected error for mg error, got nil")
	}
}

func TestExerciseService_GetExercisesByEquipment(t *testing.T) {
	repo := &mockRepository{
		GetExercisesByEquipmentFunc: func(eq []string) ([]*dto.ExerciseResponseDTO, error) {
			if len(eq) == 1 && eq[0] == "fail" {
				return nil, errors.New("eq error")
			}
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	service := NewExerciseService(repo, &mockEquipmentService{}, &mockMuscularGroupService{})
	res, err := service.GetExercisesByEquipment([]string{"eq1"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0].ID != "id1" {
		t.Errorf("expected one result with id 'id1', got %v", res)
	}
	_, err = service.GetExercisesByEquipment([]string{"fail"})
	if err == nil {
		t.Error("expected error for eq error, got nil")
	}
}

func TestExerciseService_GetExercisesByMuscularGroupAndEquipment(t *testing.T) {
	repo := &mockRepository{
		GetExercisesByMuscularGroupFunc: func(mg []string) ([]*dto.ExerciseResponseDTO, error) {
			if len(mg) == 1 && mg[0] == "fail" {
				return nil, errors.New("mg error")
			}
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
		GetExercisesByEquipmentFunc: func(eq []string) ([]*dto.ExerciseResponseDTO, error) {
			if len(eq) == 1 && eq[0] == "fail" {
				return nil, errors.New("eq error")
			}
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
		GetAllExercisesFunc: func() ([]*dto.ExerciseResponseDTO, error) {
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	service := NewExerciseService(repo, &mockEquipmentService{}, &mockMuscularGroupService{})

	// Both filters
	res, err := service.GetExercisesByMuscularGroupAndEquipment([]string{"mg1"}, []string{"eq1"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0].ID != "id1" {
		t.Errorf("expected one result with id 'id1', got %v", res)
	}

	// Only muscular group
	res, err = service.GetExercisesByMuscularGroupAndEquipment([]string{"mg1"}, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0].ID != "id1" {
		t.Errorf("expected one result with id 'id1', got %v", res)
	}

	// Only equipment
	res, err = service.GetExercisesByMuscularGroupAndEquipment(nil, []string{"eq1"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0].ID != "id1" {
		t.Errorf("expected one result with id 'id1', got %v", res)
	}

	// No filters
	res, err = service.GetExercisesByMuscularGroupAndEquipment(nil, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0].ID != "id1" {
		t.Errorf("expected one result with id 'id1', got %v", res)
	}

	// Muscular group error
	_, err = service.GetExercisesByMuscularGroupAndEquipment([]string{"fail"}, nil)
	if err == nil {
		t.Error("expected error for mg error, got nil")
	}

	// Equipment error
	_, err = service.GetExercisesByMuscularGroupAndEquipment(nil, []string{"fail"})
	if err == nil {
		t.Error("expected error for eq error, got nil")
	}
}
