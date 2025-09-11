package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	createFunc               func(gymID string, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error)
	getByIDFunc              func(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error)
	getSummaryByIDFunc       func(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error)
	getByUserIDFunc          func(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	getSummariesByUserIDFunc func(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error)
	getLastsByUserIDFunc     func(gymID, userID string, numberOfWorkouts int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	listFunc                 func(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error)
	listSummariesFunc        func(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error)
	updateFunc               func(gymID string, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error
	deleteFunc               func(gymID, id string) error
}

func (m *mockService) CreateCustomWorkoutInstance(gymID string, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error) {
	return m.createFunc(gymID, createdBy, instance)
}

func (m *mockService) GetCustomWorkoutInstanceByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
	return m.getByIDFunc(gymID, id)
}

func (m *mockService) GetCustomWorkoutInstanceSummaryByID(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error) {
	return m.getSummaryByIDFunc(gymID, id)
}

func (m *mockService) GetCustomWorkoutInstancesByUserID(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	return m.getByUserIDFunc(gymID, userID)
}

func (m *mockService) GetCustomWorkoutInstanceSummariesByUserID(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	return m.getSummariesByUserIDFunc(gymID, userID)
}

func (m *mockService) GetLastCustomWorkoutInstancesByUserID(gymID, userID string, numberOfWorkouts int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	return m.getLastsByUserIDFunc(gymID, userID, numberOfWorkouts)
}

func (m *mockService) ListCustomWorkoutInstances(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	return m.listFunc(gymID)
}

func (m *mockService) ListCustomWorkoutInstanceSummaries(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	return m.listSummariesFunc(gymID)
}

func (m *mockService) UpdateCustomWorkoutInstance(gymID string, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error {
	return m.updateFunc(gymID, id, instance)
}

func (m *mockService) DeleteCustomWorkoutInstance(gymID, id string) error {
	return m.deleteFunc(gymID, id)
}

func setupRouter(service *mockService) *chi.Mux {
	router := chi.NewRouter()
	handler := handler.NewCustomWorkoutInstanceHandler(service)

	// Add URL parameters for testing
	router.Route("/gym/{gymID}", func(r chi.Router) {
		r.Post("/custom-workout-instance", handler.Create)
		r.Get("/custom-workout-instance/{id}", handler.GetByID)
		r.Get("/custom-workout-instance/{id}/summary", handler.GetSummaryByID)
		r.Get("/custom-workout-instance/user/{userID}", handler.GetByUserID)
		r.Get("/custom-workout-instance/user/{userID}/summaries", handler.GetSummariesByUserID)
		r.Get("/custom-workout-instance/user/{userID}/last", handler.GetLastsByUserID)
		r.Get("/custom-workout-instance", handler.List)
		r.Get("/custom-workout-instance/summaries", handler.ListSummaries)
		r.Put("/custom-workout-instance/{id}", handler.Update)
		r.Delete("/custom-workout-instance/{id}", handler.Delete)
	})

	return router
}

func TestCreate_Success(t *testing.T) {
	createdID := "instance123"
	service := &mockService{
		createFunc: func(gymID string, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "user123", createdBy)
			assert.Equal(t, "Test Workout", instance.Name)
			return &createdID, nil
		},
	}

	router := setupRouter(service)

	body := map[string]interface{}{
		"name":            "Test Workout",
		"description":     "Test description",
		"template_source": "gym",
		"gym_template_id": "template123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/gym/gym123/custom-workout-instance", bytes.NewReader(bodyBytes))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user123"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreate_InvalidJSON(t *testing.T) {
	service := &mockService{}
	router := setupRouter(service)

	req := httptest.NewRequest("POST", "/gym/gym123/custom-workout-instance", bytes.NewReader([]byte("invalid json")))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user123"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_ServiceError(t *testing.T) {
	service := &mockService{
		createFunc: func(gymID string, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error) {
			return nil, apierror.New(errorcode_enum.CodeBadRequest, "Validation error", nil)
		},
	}

	router := setupRouter(service)

	body := map[string]interface{}{
		"name":            "Test Workout",
		"template_source": "gym",
		"gym_template_id": "template123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/gym/gym123/custom-workout-instance", bytes.NewReader(bodyBytes))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user123"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByID_Success(t *testing.T) {
	instance := &dto.ResponseCustomWorkoutInstanceDTO{
		ID:             "instance123",
		Name:           "Test Workout",
		TemplateSource: "gym",
	}

	service := &mockService{
		getByIDFunc: func(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "instance123", id)
			return instance, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/instance123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetByID_NotFound(t *testing.T) {
	service := &mockService{
		getByIDFunc: func(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", nil)
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSummaryByID_Success(t *testing.T) {
	summary := &dto.SummaryCustomWorkoutInstanceDTO{
		ID:             "instance123",
		Name:           "Test Workout",
		TemplateSource: "gym",
	}

	service := &mockService{
		getSummaryByIDFunc: func(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "instance123", id)
			return summary, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/instance123/summary", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetByUserID_Success(t *testing.T) {
	instances := []*dto.ResponseCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Workout 1",
			TemplateSource: "gym",
		},
		{
			ID:             "instance2",
			Name:           "Workout 2",
			TemplateSource: "public",
		},
	}

	service := &mockService{
		getByUserIDFunc: func(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "user123", userID)
			return instances, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/user/user123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSummariesByUserID_Success(t *testing.T) {
	summaries := []*dto.SummaryCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Workout 1",
			TemplateSource: "gym",
		},
	}

	service := &mockService{
		getSummariesByUserIDFunc: func(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "user123", userID)
			return summaries, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/user/user123/summaries", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetLastsByUserID_Success(t *testing.T) {
	instances := []*dto.ResponseCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Recent Workout",
			TemplateSource: "gym",
		},
	}

	service := &mockService{
		getLastsByUserIDFunc: func(gymID, userID string, count int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "user123", userID)
			assert.Equal(t, 5, count) // Default count
			return instances, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/user/user123/last?count=5", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetLastsByUserID_CustomCount(t *testing.T) {
	instances := []*dto.ResponseCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Recent Workout",
			TemplateSource: "gym",
		},
	}

	service := &mockService{
		getLastsByUserIDFunc: func(gymID, userID string, count int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "user123", userID)
			assert.Equal(t, 3, count) // Custom count from query param
			return instances, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/user/user123/last?count=3", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestList_Success(t *testing.T) {
	instances := []*dto.ResponseCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Workout 1",
			TemplateSource: "gym",
		},
		{
			ID:             "instance2",
			Name:           "Workout 2",
			TemplateSource: "public",
		},
	}

	service := &mockService{
		listFunc: func(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			return instances, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListSummaries_Success(t *testing.T) {
	summaries := []*dto.SummaryCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Workout 1",
			TemplateSource: "gym",
		},
		{
			ID:             "instance2",
			Name:           "Workout 2",
			TemplateSource: "public",
		},
	}

	service := &mockService{
		listSummariesFunc: func(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
			assert.Equal(t, "gym123", gymID)
			return summaries, nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("GET", "/gym/gym123/custom-workout-instance/summaries", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	service := &mockService{
		updateFunc: func(gymID string, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "instance123", id)
			assert.Equal(t, "Updated Workout", *instance.Name)
			return nil
		},
	}

	router := setupRouter(service)

	body := map[string]interface{}{
		"name":        "Updated Workout",
		"description": "Updated description",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/gym/gym123/custom-workout-instance/instance123", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate_InvalidJSON(t *testing.T) {
	service := &mockService{}
	router := setupRouter(service)

	req := httptest.NewRequest("PUT", "/gym/gym123/custom-workout-instance/instance123", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_ServiceError(t *testing.T) {
	service := &mockService{
		updateFunc: func(gymID string, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error {
			return apierror.New(errorcode_enum.CodeBadRequest, "Validation error", nil)
		},
	}

	router := setupRouter(service)

	body := map[string]interface{}{
		"name": "Updated Workout",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/gym/gym123/custom-workout-instance/instance123", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDelete_Success(t *testing.T) {
	service := &mockService{
		deleteFunc: func(gymID, id string) error {
			assert.Equal(t, "gym123", gymID)
			assert.Equal(t, "instance123", id)
			return nil
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("DELETE", "/gym/gym123/custom-workout-instance/instance123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDelete_ServiceError(t *testing.T) {
	service := &mockService{
		deleteFunc: func(gymID, id string) error {
			return apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", nil)
		},
	}

	router := setupRouter(service)

	req := httptest.NewRequest("DELETE", "/gym/gym123/custom-workout-instance/instance123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

