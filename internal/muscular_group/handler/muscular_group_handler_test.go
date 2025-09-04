package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	createFunc func(*dto.CreateMuscularGroupDTO) (string, error)
	getFunc    func(string) (*dto.MuscularGroupResponseDTO, error)
	listFunc   func() ([]*dto.MuscularGroupResponseDTO, error)
	updateFunc func(string, *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error)
	deleteFunc func(string) error
}

func (m *mockService) CreateMuscularGroup(dto *dto.CreateMuscularGroupDTO) (string, error) {
	return m.createFunc(dto)
}
func (m *mockService) GetMuscularGroupByID(id string) (*dto.MuscularGroupResponseDTO, error) {
	return m.getFunc(id)
}
func (m *mockService) GetAllMuscularGroups() ([]*dto.MuscularGroupResponseDTO, error) {
	return m.listFunc()
}
func (m *mockService) UpdateMuscularGroup(id string, dto *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error) {
	return m.updateFunc(id, dto)
}
func (m *mockService) DeleteMuscularGroup(id string) error {
	return m.deleteFunc(id)
}

func TestCreateMuscularGroupHandler(t *testing.T) {
	service := &mockService{
		createFunc: func(dto *dto.CreateMuscularGroupDTO) (string, error) {
			if dto.Name == "error" {
				return "", apierror.New(enum.CodeConflict, "Duplicate", nil)
			}
			return "id123", nil
		},
	}
	h := handler.NewMuscularGroupHandler(service)
	body, _ := json.Marshal(dto.CreateMuscularGroupDTO{Name: "Chest"})
	req := httptest.NewRequest(http.MethodPost, "/muscular-groups", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateMuscularGroup(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Invalid body
	req = httptest.NewRequest(http.MethodPost, "/muscular-groups", bytes.NewReader([]byte("bad json")))
	w = httptest.NewRecorder()
	h.CreateMuscularGroup(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Service error
	body, _ = json.Marshal(dto.CreateMuscularGroupDTO{Name: "error"})
	req = httptest.NewRequest(http.MethodPost, "/muscular-groups", bytes.NewReader(body))
	w = httptest.NewRecorder()
	h.CreateMuscularGroup(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetMuscularGroupHandler(t *testing.T) {
	service := &mockService{
		getFunc: func(id string) (*dto.MuscularGroupResponseDTO, error) {
			if id == "notfound" {
				return nil, apierror.New(enum.CodeNotFound, "Not found", nil)
			}
			return &dto.MuscularGroupResponseDTO{ID: id, Name: "Chest"}, nil
		},
	}
	h := handler.NewMuscularGroupHandler(service)
	r := chi.NewRouter()
	r.Get("/muscular-groups/{id}", h.GetMuscularGroup)

	req := httptest.NewRequest(http.MethodGet, "/muscular-groups/id123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Missing ID (should be 404 because route doesn't match)
	req = httptest.NewRequest(http.MethodGet, "/muscular-groups/", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Not found
	req = httptest.NewRequest(http.MethodGet, "/muscular-groups/notfound", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListMuscularGroupsHandler(t *testing.T) {
	service := &mockService{
		listFunc: func() ([]*dto.MuscularGroupResponseDTO, error) {
			return []*dto.MuscularGroupResponseDTO{{ID: "1", Name: "Chest"}}, nil
		},
	}
	h := handler.NewMuscularGroupHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/muscular-groups", nil)
	w := httptest.NewRecorder()
	h.ListMuscularGroups(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Service error
	service.listFunc = func() ([]*dto.MuscularGroupResponseDTO, error) {
		return nil, apierror.New(enum.CodeInternal, "fail", nil)
	}
	req = httptest.NewRequest(http.MethodGet, "/muscular-groups", nil)
	w = httptest.NewRecorder()
	h.ListMuscularGroups(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateMuscularGroupHandler(t *testing.T) {
	service := &mockService{
		updateFunc: func(id string, muscularGroupDto *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error) {
			if id == "notfound" {
				return nil, apierror.New(enum.CodeNotFound, "Not found", nil)
			}
			if muscularGroupDto.Name != nil && *muscularGroupDto.Name == "conflict" {
				return nil, apierror.New(enum.CodeConflict, "Conflict", nil)
			}
			if muscularGroupDto.Name != nil && *muscularGroupDto.Name == "fail" {
				return nil, apierror.New(enum.CodeInternal, "fail", nil)
			}
			return &dto.MuscularGroupResponseDTO{ID: id, Name: "Updated"}, nil
		},
	}
	h := handler.NewMuscularGroupHandler(service)
	r := chi.NewRouter()
	r.Put("/muscular-groups/{id}", h.UpdateMuscularGroup)

	name := "Updated"
	body, _ := json.Marshal(dto.UpdateMuscularGroupDTO{Name: &name})
	req := httptest.NewRequest(http.MethodPut, "/muscular-groups/id123", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Missing ID (should be 404 because route doesn't match)
	req = httptest.NewRequest(http.MethodPut, "/muscular-groups/", bytes.NewReader(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Invalid body
	req = httptest.NewRequest(http.MethodPut, "/muscular-groups/id123", bytes.NewReader([]byte("bad json")))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Not found
	name = "notfound"
	body, _ = json.Marshal(dto.UpdateMuscularGroupDTO{Name: &name})
	req = httptest.NewRequest(http.MethodPut, "/muscular-groups/notfound", bytes.NewReader(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Conflict
	name = "conflict"
	body, _ = json.Marshal(dto.UpdateMuscularGroupDTO{Name: &name})
	req = httptest.NewRequest(http.MethodPut, "/muscular-groups/id123", bytes.NewReader(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)

	// Internal error
	name = "fail"
	body, _ = json.Marshal(dto.UpdateMuscularGroupDTO{Name: &name})
	req = httptest.NewRequest(http.MethodPut, "/muscular-groups/id123", bytes.NewReader(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteMuscularGroupHandler(t *testing.T) {
	service := &mockService{
		deleteFunc: func(id string) error {
			if id == "notfound" {
				return apierror.New(enum.CodeNotFound, "Not found", nil)
			}
			if id == "fail" {
				return apierror.New(enum.CodeInternal, "fail", nil)
			}
			return nil
		},
	}
	h := handler.NewMuscularGroupHandler(service)
	r := chi.NewRouter()
	r.Delete("/muscular-groups/{id}", h.DeleteMuscularGroup)

	req := httptest.NewRequest(http.MethodDelete, "/muscular-groups/id123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Missing ID (should be 404 because route doesn't match)
	req = httptest.NewRequest(http.MethodDelete, "/muscular-groups/", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Not found
	req = httptest.NewRequest(http.MethodDelete, "/muscular-groups/notfound", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Internal error
	req = httptest.NewRequest(http.MethodDelete, "/muscular-groups/fail", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
