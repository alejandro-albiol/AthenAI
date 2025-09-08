package service

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	CreateFn  func(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error)
	GetByIDFn func(gymID, id string) (*dto.CustomExerciseResponseDTO, error)
	ListFn    func(gymID string) ([]*dto.CustomExerciseResponseDTO, error)
	DeleteFn  func(gymID, id string) error
	UpdateFn  func(gymID, id string, update *dto.CustomExerciseUpdateDTO) error
}

func (m *mockRepo) UpdateCustomExercise(gymID, id string, update *dto.CustomExerciseUpdateDTO) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(gymID, id, update)
	}
	return nil
}

func (m *mockRepo) CreateCustomExercise(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error) {
	return m.CreateFn(gymID, exercise)
}
func (m *mockRepo) GetCustomExerciseByID(gymID, id string) (*dto.CustomExerciseResponseDTO, error) {
	return m.GetByIDFn(gymID, id)
}
func (m *mockRepo) ListCustomExercises(gymID string) ([]*dto.CustomExerciseResponseDTO, error) {
	return m.ListFn(gymID)
}
func (m *mockRepo) DeleteCustomExercise(gymID, id string) error {
	return m.DeleteFn(gymID, id)
}

func TestCreateCustomExerciseService(t *testing.T) {
	repo := &mockRepo{
		CreateFn: func(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error) {
			if exercise.Name == "" {
				return nil, errors.New("name required")
			}
			id := "ex-1"
			return &id, nil
		},
	}
	svc := NewCustomExerciseService(repo)
	dtoReq := &dto.CustomExerciseCreationDTO{
		Name:            "Push Up",
		Synonyms:        []string{"Pushup"},
		DifficultyLevel: "easy",
		ExerciseType:    "bodyweight",
		Instructions:    "Do a push up",
		VideoURL:        "http://example.com/video",
		ImageURL:        "http://example.com/image",
		MuscularGroups:  []string{"chest", "triceps"},
		CreatedBy:       "user1",
	}
	id, err := svc.CreateCustomExercise("tenant1", dtoReq)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "ex-1", *id)
}

func TestGetCustomExerciseByIDService(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(gymID, id string) (*dto.CustomExerciseResponseDTO, error) {
			if id == "notfound" {
				return nil, errors.New("not found")
			}
			return &dto.CustomExerciseResponseDTO{ID: id, Name: "Push Up"}, nil
		},
	}
	svc := NewCustomExerciseService(repo)
	result, err := svc.GetCustomExerciseByID("tenant1", "ex-1")
	assert.NoError(t, err)
	assert.Equal(t, "ex-1", result.ID)

	_, err = svc.GetCustomExerciseByID("tenant1", "notfound")
	assert.Error(t, err)
}

func TestListCustomExercisesService(t *testing.T) {
	repo := &mockRepo{
		ListFn: func(gymID string) ([]*dto.CustomExerciseResponseDTO, error) {
			return []*dto.CustomExerciseResponseDTO{{ID: "ex-1", Name: "Push Up"}}, nil
		},
	}
	svc := NewCustomExerciseService(repo)
	results, err := svc.ListCustomExercises("tenant1")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "ex-1", results[0].ID)
}

func TestDeleteCustomExerciseService(t *testing.T) {
	repo := &mockRepo{
		DeleteFn: func(gymID, id string) error {
			if id == "notfound" {
				return errors.New("not found")
			}
			return nil
		},
	}
	svc := NewCustomExerciseService(repo)
	err := svc.DeleteCustomExercise("tenant1", "ex-1")
	assert.NoError(t, err)

	err = svc.DeleteCustomExercise("tenant1", "notfound")
	assert.Error(t, err)
}
