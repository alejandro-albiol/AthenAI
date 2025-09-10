package service_test

import (
	"database/sql"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	createErr                  error
	getByIDErr                 error
	listByWorkoutInstanceIDErr error
	listByMuscularGroupIDErr   error
	listByEquipmentIDErr       error
	updateErr                  error
	deleteErr                  error
	exercises                  []*dto.ResponseCustomWorkoutExerciseDTO
	lastCreatedID              string
}

func (m *mockRepository) Create(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &m.lastCreatedID, nil
}

func (m *mockRepository) GetByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	for _, e := range m.exercises {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) ListByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if m.listByWorkoutInstanceIDErr != nil {
		return nil, m.listByWorkoutInstanceIDErr
	}
	var result []*dto.ResponseCustomWorkoutExerciseDTO
	for _, e := range m.exercises {
		if e.WorkoutInstanceID == workoutInstanceID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockRepository) ListByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if m.listByMuscularGroupIDErr != nil {
		return nil, m.listByMuscularGroupIDErr
	}
	return m.exercises, nil
}

func (m *mockRepository) ListByEquipmentID(gymID, equipmentID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if m.listByEquipmentIDErr != nil {
		return nil, m.listByEquipmentIDErr
	}
	return m.exercises, nil
}

func (m *mockRepository) Update(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error {
	return m.updateErr
}

func (m *mockRepository) Delete(gymID, id string) error {
	return m.deleteErr
}

func TestCreateCustomWorkoutExercise_Success(t *testing.T) {
	mockRepo := &mockRepository{
		lastCreatedID: "exercise123",
		exercises:     []*dto.ResponseCustomWorkoutExerciseDTO{}, // Empty list for duplicate check
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	exercise := &dto.CreateCustomWorkoutExerciseDTO{
		CreatedBy:         "user123",
		WorkoutInstanceID: "workout456",
		ExerciseSource:    "public",
		PublicExerciseID:  stringPtr("exercise789"),
		BlockName:         "main",
		ExerciseOrder:     1,
		Sets:              intPtr(3),
		RepsMin:           intPtr(8),
		RepsMax:           intPtr(12),
		WeightKg:          floatPtr(50.5),
		RestSeconds:       intPtr(60),
		Notes:             stringPtr("Test exercise"),
	}

	id, err := svc.CreateCustomWorkoutExercise(gymID, exercise)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "exercise123", *id)
}

func TestCreateCustomWorkoutExercise_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		exercise    *dto.CreateCustomWorkoutExerciseDTO
		expectedErr string
	}{
		{
			name: "Empty CreatedBy",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "public",
				PublicExerciseID:  stringPtr("exercise789"),
				BlockName:         "main",
				ExerciseOrder:     1,
			},
			expectedErr: "CreatedBy is required",
		},
		{
			name: "Empty WorkoutInstanceID",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:        "user123",
				ExerciseSource:   "public",
				PublicExerciseID: stringPtr("exercise789"),
				BlockName:        "main",
				ExerciseOrder:    1,
			},
			expectedErr: "WorkoutInstanceID is required",
		},
		{
			name: "Invalid ExerciseSource",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:         "user123",
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "invalid",
				BlockName:         "main",
				ExerciseOrder:     1,
			},
			expectedErr: "ExerciseSource must be 'public' or 'gym'",
		},
		{
			name: "Public source without PublicExerciseID",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:         "user123",
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "public",
				BlockName:         "main",
				ExerciseOrder:     1,
			},
			expectedErr: "PublicExerciseID is required when ExerciseSource is 'public'",
		},
		{
			name: "Gym source without GymExerciseID",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:         "user123",
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "gym",
				BlockName:         "main",
				ExerciseOrder:     1,
			},
			expectedErr: "GymExerciseID is required when ExerciseSource is 'gym'",
		},
		{
			name: "Invalid ExerciseOrder",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:         "user123",
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "public",
				PublicExerciseID:  stringPtr("exercise789"),
				BlockName:         "main",
				ExerciseOrder:     0,
			},
			expectedErr: "ExerciseOrder must be greater than 0",
		},
		{
			name: "Invalid Sets",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:         "user123",
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "public",
				PublicExerciseID:  stringPtr("exercise789"),
				BlockName:         "main",
				ExerciseOrder:     1,
				Sets:              intPtr(-1),
			},
			expectedErr: "Sets must be greater than 0",
		},
		{
			name: "RepsMin greater than RepsMax",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:         "user123",
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "public",
				PublicExerciseID:  stringPtr("exercise789"),
				BlockName:         "main",
				ExerciseOrder:     1,
				RepsMin:           intPtr(15),
				RepsMax:           intPtr(10),
			},
			expectedErr: "RepsMin cannot be greater than RepsMax",
		},
		{
			name: "Negative WeightKg",
			exercise: &dto.CreateCustomWorkoutExerciseDTO{
				CreatedBy:         "user123",
				WorkoutInstanceID: "workout456",
				ExerciseSource:    "public",
				PublicExerciseID:  stringPtr("exercise789"),
				BlockName:         "main",
				ExerciseOrder:     1,
				WeightKg:          floatPtr(-10.5),
			},
			expectedErr: "WeightKg cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{exercises: []*dto.ResponseCustomWorkoutExerciseDTO{}}
			svc := service.NewCustomWorkoutExerciseService(mockRepo)
			gymID := "gym123"

			id, err := svc.CreateCustomWorkoutExercise(gymID, tt.exercise)

			assert.Error(t, err)
			assert.Nil(t, id)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestCreateCustomWorkoutExercise_DuplicateOrderConflict(t *testing.T) {
	mockRepo := &mockRepository{
		exercises: []*dto.ResponseCustomWorkoutExerciseDTO{
			{
				ID:                "existing123",
				WorkoutInstanceID: "workout456",
				BlockName:         "main",
				ExerciseOrder:     1,
			},
		},
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	exercise := &dto.CreateCustomWorkoutExerciseDTO{
		CreatedBy:         "user123",
		WorkoutInstanceID: "workout456",
		ExerciseSource:    "public",
		PublicExerciseID:  stringPtr("exercise789"),
		BlockName:         "main",
		ExerciseOrder:     1, // Duplicate order
	}

	id, err := svc.CreateCustomWorkoutExercise(gymID, exercise)

	assert.Error(t, err)
	assert.Nil(t, id)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeConflict, apiErr.Code)
	assert.Contains(t, err.Error(), "Exercise order 1 already exists in block 'main'")
}

func TestGetCustomWorkoutExerciseByID_Success(t *testing.T) {
	expectedExercise := &dto.ResponseCustomWorkoutExerciseDTO{
		ID:                "exercise123",
		CreatedBy:         "user123",
		WorkoutInstanceID: "workout456",
		ExerciseSource:    "public",
		PublicExerciseID:  stringPtr("exercise789"),
		BlockName:         "main",
		ExerciseOrder:     1,
	}

	mockRepo := &mockRepository{
		exercises: []*dto.ResponseCustomWorkoutExerciseDTO{expectedExercise},
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	result, err := svc.GetCustomWorkoutExerciseByID(gymID, "exercise123")

	assert.NoError(t, err)
	assert.Equal(t, expectedExercise, result)
}

func TestGetCustomWorkoutExerciseByID_NotFound(t *testing.T) {
	mockRepo := &mockRepository{
		getByIDErr: sql.ErrNoRows,
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	result, err := svc.GetCustomWorkoutExerciseByID(gymID, "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
}

func TestListCustomWorkoutExercisesByWorkoutInstanceID_Success(t *testing.T) {
	expectedExercises := []*dto.ResponseCustomWorkoutExerciseDTO{
		{
			ID:                "ex1",
			WorkoutInstanceID: "workout456",
			BlockName:         "warmup",
			ExerciseOrder:     1,
		},
		{
			ID:                "ex2",
			WorkoutInstanceID: "workout456",
			BlockName:         "main",
			ExerciseOrder:     1,
		},
	}

	mockRepo := &mockRepository{
		exercises: expectedExercises,
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	result, err := svc.ListCustomWorkoutExercisesByWorkoutInstanceID(gymID, "workout456")

	assert.NoError(t, err)
	assert.Equal(t, expectedExercises, result)
}

func TestListCustomWorkoutExercisesByMuscularGroupID_Success(t *testing.T) {
	expectedExercises := []*dto.ResponseCustomWorkoutExerciseDTO{
		{ID: "ex1", Notes: stringPtr("Chest exercise")},
	}

	mockRepo := &mockRepository{
		exercises: expectedExercises,
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	result, err := svc.ListCustomWorkoutExercisesByMuscularGroupID(gymID, "muscle123")

	assert.NoError(t, err)
	assert.Equal(t, expectedExercises, result)
}

func TestListCustomWorkoutExercisesByEquipmentID_Success(t *testing.T) {
	expectedExercises := []*dto.ResponseCustomWorkoutExerciseDTO{
		{ID: "ex1", Notes: stringPtr("Barbell exercise")},
	}

	mockRepo := &mockRepository{
		exercises: expectedExercises,
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	result, err := svc.ListCustomWorkoutExercisesByEquipmentID(gymID, "equipment123")

	assert.NoError(t, err)
	assert.Equal(t, expectedExercises, result)
}

func TestUpdateCustomWorkoutExercise_Success(t *testing.T) {
	existingExercise := &dto.ResponseCustomWorkoutExerciseDTO{
		ID: "exercise123",
	}

	mockRepo := &mockRepository{
		exercises: []*dto.ResponseCustomWorkoutExerciseDTO{existingExercise},
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	updateDTO := &dto.UpdateCustomWorkoutExerciseDTO{
		ID:      "exercise123",
		Sets:    intPtr(4),
		RepsMin: intPtr(6),
		RepsMax: intPtr(10),
	}

	err := svc.UpdateCustomWorkoutExercise(gymID, updateDTO)

	assert.NoError(t, err)
}

func TestUpdateCustomWorkoutExercise_NotFound(t *testing.T) {
	mockRepo := &mockRepository{
		getByIDErr: sql.ErrNoRows,
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	updateDTO := &dto.UpdateCustomWorkoutExerciseDTO{
		ID: "nonexistent",
	}

	err := svc.UpdateCustomWorkoutExercise(gymID, updateDTO)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
}

func TestDeleteCustomWorkoutExercise_Success(t *testing.T) {
	existingExercise := &dto.ResponseCustomWorkoutExerciseDTO{
		ID: "exercise123",
	}

	mockRepo := &mockRepository{
		exercises: []*dto.ResponseCustomWorkoutExerciseDTO{existingExercise},
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	err := svc.DeleteCustomWorkoutExercise(gymID, "exercise123")

	assert.NoError(t, err)
}

func TestDeleteCustomWorkoutExercise_NotFound(t *testing.T) {
	mockRepo := &mockRepository{
		getByIDErr: sql.ErrNoRows,
	}
	svc := service.NewCustomWorkoutExerciseService(mockRepo)
	gymID := "gym123"

	err := svc.DeleteCustomWorkoutExercise(gymID, "nonexistent")

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}
