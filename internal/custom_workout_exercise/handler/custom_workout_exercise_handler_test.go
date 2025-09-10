package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	createFunc                  func(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error)
	getByIDFunc                 func(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error)
	listByWorkoutInstanceIDFunc func(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error)
	listByMuscularGroupIDFunc   func(gymID, muscularGroupID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error)
	listByEquipmentIDFunc       func(gymID, equipmentID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error)
	updateFunc                  func(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error
	deleteFunc                  func(gymID, id string) error
}

func (m *mockService) CreateCustomWorkoutExercise(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error) {
	return m.createFunc(gymID, exercise)
}

func (m *mockService) GetCustomWorkoutExerciseByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
	return m.getByIDFunc(gymID, id)
}

func (m *mockService) ListCustomWorkoutExercisesByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	return m.listByWorkoutInstanceIDFunc(gymID, workoutInstanceID)
}

func (m *mockService) ListCustomWorkoutExercisesByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	return m.listByMuscularGroupIDFunc(gymID, muscularGroupID)
}

func (m *mockService) ListCustomWorkoutExercisesByEquipmentID(gymID, equipmentID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	return m.listByEquipmentIDFunc(gymID, equipmentID)
}

func (m *mockService) UpdateCustomWorkoutExercise(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error {
	return m.updateFunc(gymID, exercise)
}

func (m *mockService) DeleteCustomWorkoutExercise(gymID, id string) error {
	return m.deleteFunc(gymID, id)
}

func TestCreateCustomWorkoutExerciseHandler_Success(t *testing.T) {
	service := &mockService{
		createFunc: func(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "user123", exercise.CreatedBy)
			id := "exercise123"
			return &id, nil
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	exercise := dto.CreateCustomWorkoutExerciseDTO{
		CreatedBy:         "user123",
		WorkoutInstanceID: "workout456",
		ExerciseSource:    "public",
		PublicExerciseID:  stringPtr("exercise789"),
		BlockName:         "main",
		ExerciseOrder:     1,
		Sets:              intPtr(3),
		RepsMin:           intPtr(8),
		RepsMax:           intPtr(12),
	}

	body, _ := json.Marshal(exercise)
	req := httptest.NewRequest(http.MethodPost, "/custom-workout-exercises", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Mock middleware values
	ctx := req.Context()
	ctx = contextWithGymID(ctx, "gym123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Workout exercise created successfully", response["message"])
}

func TestCreateCustomWorkoutExerciseHandler_InvalidJSON(t *testing.T) {
	service := &mockService{}
	h := handler.NewCustomWorkoutExerciseHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/custom-workout-exercises", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	ctx := contextWithGymID(req.Context(), "gym123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateCustomWorkoutExerciseHandler_ServiceError(t *testing.T) {
	service := &mockService{
		createFunc: func(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error) {
			return nil, apierror.New(errorcode_enum.CodeBadRequest, "Validation error", nil)
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	exercise := dto.CreateCustomWorkoutExerciseDTO{
		CreatedBy:      "user123",
		ExerciseSource: "public",
		BlockName:      "main",
		ExerciseOrder:  1,
	}

	body, _ := json.Marshal(exercise)
	req := httptest.NewRequest(http.MethodPost, "/custom-workout-exercises", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := contextWithGymID(req.Context(), "gym123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByIDHandler_Success(t *testing.T) {
	expectedExercise := &dto.ResponseCustomWorkoutExerciseDTO{
		ID:             "exercise123",
		CreatedBy:      "user123",
		ExerciseSource: "public",
		BlockName:      "main",
		ExerciseOrder:  1,
	}

	service := &mockService{
		getByIDFunc: func(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "exercise123", id)
			return expectedExercise, nil
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/custom-workout-exercises/exercise123", nil)

	ctx := contextWithGymID(req.Context(), "gym123")
	ctx = contextWithURLParam(ctx, "id", "exercise123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Workout exercise retrieved successfully", response["message"])
}

func TestGetByIDHandler_NotFound(t *testing.T) {
	service := &mockService{
		getByIDFunc: func(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Exercise not found", nil)
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/custom-workout-exercises/nonexistent", nil)

	ctx := contextWithGymID(req.Context(), "gym123")
	ctx = contextWithURLParam(ctx, "id", "nonexistent")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListByWorkoutInstanceIDHandler_Success(t *testing.T) {
	expectedExercises := []*dto.ResponseCustomWorkoutExerciseDTO{
		{ID: "ex1", BlockName: "warmup"},
		{ID: "ex2", BlockName: "main"},
	}

	service := &mockService{
		listByWorkoutInstanceIDFunc: func(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "workout456", workoutInstanceID)
			return expectedExercises, nil
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/custom-workout-exercises/workout-instance/workout456", nil)

	ctx := contextWithGymID(req.Context(), "gym123")
	ctx = contextWithURLParam(ctx, "workoutInstanceId", "workout456")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.ListByWorkoutInstanceID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Workout exercises retrieved successfully", response["message"])
}

func TestListByMuscularGroupIDHandler_Success(t *testing.T) {
	expectedExercises := []*dto.ResponseCustomWorkoutExerciseDTO{
		{ID: "ex1", Notes: stringPtr("Chest exercise")},
	}

	service := &mockService{
		listByMuscularGroupIDFunc: func(gymID, muscularGroupID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "muscle123", muscularGroupID)
			return expectedExercises, nil
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/custom-workout-exercises/muscular-group/muscle123", nil)

	ctx := contextWithGymID(req.Context(), "gym123")
	ctx = contextWithURLParam(ctx, "muscularGroupId", "muscle123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.ListByMuscularGroupID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListByEquipmentIDHandler_Success(t *testing.T) {
	expectedExercises := []*dto.ResponseCustomWorkoutExerciseDTO{
		{ID: "ex1", Notes: stringPtr("Barbell exercise")},
	}

	service := &mockService{
		listByEquipmentIDFunc: func(gymID, equipmentID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "equipment123", equipmentID)
			return expectedExercises, nil
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/custom-workout-exercises/equipment/equipment123", nil)

	ctx := contextWithGymID(req.Context(), "gym123")
	ctx = contextWithURLParam(ctx, "equipmentId", "equipment123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.ListByEquipmentID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateHandler_Success(t *testing.T) {
	service := &mockService{
		updateFunc: func(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "exercise123", exercise.ID)
			assert.Equal(t, 4, *exercise.Sets)
			return nil
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	updateDTO := dto.UpdateCustomWorkoutExerciseDTO{
		Sets:    intPtr(4),
		RepsMin: intPtr(6),
		RepsMax: intPtr(10),
	}

	body, _ := json.Marshal(updateDTO)
	req := httptest.NewRequest(http.MethodPut, "/custom-workout-exercises/exercise123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := contextWithGymID(req.Context(), "gym123")
	ctx = contextWithURLParam(ctx, "id", "exercise123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteHandler_Success(t *testing.T) {
	service := &mockService{
		deleteFunc: func(gymID, id string) error {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "exercise123", id)
			return nil
		},
	}

	h := handler.NewCustomWorkoutExerciseHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/custom-workout-exercises/exercise123", nil)

	ctx := contextWithGymID(req.Context(), "gym123")
	ctx = contextWithURLParam(ctx, "id", "exercise123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Delete(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func contextWithGymID(ctx context.Context, gymID string) context.Context {
	// Use the same context key as the middleware
	return context.WithValue(ctx, middleware.GymIDKey, gymID)
}

func contextWithURLParam(ctx context.Context, key, value string) context.Context {
	// Simulate chi URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return context.WithValue(ctx, chi.RouteCtxKey, rctx)
}
