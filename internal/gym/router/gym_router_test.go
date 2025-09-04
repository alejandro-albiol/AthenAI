package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/gym/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGymHandler struct {
	mock.Mock
}

func (m *MockGymHandler) CreateGym(w http.ResponseWriter, r *http.Request)    { m.Called(w, r) }
func (m *MockGymHandler) GetGymByID(w http.ResponseWriter, r *http.Request)   { m.Called(w, r) }
func (m *MockGymHandler) GetGymByName(w http.ResponseWriter, r *http.Request) { m.Called(w, r) }
func (m *MockGymHandler) GetAllGyms(w http.ResponseWriter, r *http.Request)   { m.Called(w, r) }
func (m *MockGymHandler) UpdateGym(w http.ResponseWriter, r *http.Request)    { m.Called(w, r) }
func (m *MockGymHandler) SetGymActive(w http.ResponseWriter, r *http.Request) { m.Called(w, r) }
func (m *MockGymHandler) DeleteGym(w http.ResponseWriter, r *http.Request)    { m.Called(w, r) }

func TestGymRouter_Routes(t *testing.T) {
	routes := []struct {
		name           string
		method         string
		path           string
		handler        string
		expectedStatus int
	}{
		{"get all gyms", http.MethodGet, "/", "GetAllGyms", http.StatusOK},
		{"create gym", http.MethodPost, "/", "CreateGym", http.StatusCreated},
		{"get gym by id", http.MethodGet, "/gym123", "GetGymByID", http.StatusOK},
		{"get gym by name", http.MethodGet, "/name/test-gym", "GetGymByName", http.StatusOK},
		{"update gym", http.MethodPut, "/gym123", "UpdateGym", http.StatusOK},
		{"activate gym", http.MethodPut, "/gym123/activate", "SetGymActive", http.StatusOK},
		{"deactivate gym", http.MethodPut, "/gym123/deactivate", "SetGymActive", http.StatusOK},
		{"delete gym", http.MethodDelete, "/gym123", "DeleteGym", http.StatusNoContent},
	}

	mockHandler := new(MockGymHandler)
	r := router.NewGymRouter(mockHandler)

	for _, route := range routes {
		t.Run(route.name, func(t *testing.T) {
			mockHandler.On(route.handler, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(route.expectedStatus)
			}).Once()

			req := httptest.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, route.expectedStatus, w.Code)
			mockHandler.AssertExpectations(t)
		})
	}

	t.Run("unknown route returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unknown/fakerouter/invent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("unsupported method returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/gym123", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
