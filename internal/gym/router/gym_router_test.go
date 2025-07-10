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

func (m *MockGymHandler) CreateGym(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockGymHandler) GetGymByID(w http.ResponseWriter, r *http.Request, id string) {
	m.Called(w, r, id)
}

func (m *MockGymHandler) GetGymByDomain(w http.ResponseWriter, r *http.Request, domain string) {
	m.Called(w, r, domain)
}

func (m *MockGymHandler) GetAllGyms(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockGymHandler) UpdateGym(w http.ResponseWriter, r *http.Request, id string) {
	m.Called(w, r, id)
}

func (m *MockGymHandler) SetGymActive(w http.ResponseWriter, r *http.Request, id string, active bool) {
	m.Called(w, r, id, active)
}

func (m *MockGymHandler) DeleteGym(w http.ResponseWriter, r *http.Request, id string) {
	m.Called(w, r, id)
}

func TestGymRouter(t *testing.T) {
	mockHandler := new(MockGymHandler)
	router := router.NewGymRouter(mockHandler)

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedCalls  func(*MockGymHandler)
		expectedStatus int
	}{
		{
			name:   "create gym",
			method: http.MethodPost,
			path:   "/",
			expectedCalls: func(m *MockGymHandler) {
				m.On("CreateGym", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusCreated)
				})
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "get gym by id",
			method: http.MethodGet,
			path:   "/gym123",
			expectedCalls: func(m *MockGymHandler) {
				m.On("GetGymByID", mock.Anything, mock.Anything, "gym123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "get gym by domain",
			method: http.MethodGet,
			path:   "/domain/test-gym",
			expectedCalls: func(m *MockGymHandler) {
				m.On("GetGymByDomain", mock.Anything, mock.Anything, "test-gym").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "get all gyms",
			method: http.MethodGet,
			path:   "/",
			expectedCalls: func(m *MockGymHandler) {
				m.On("GetAllGyms", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "update gym",
			method: http.MethodPut,
			path:   "/gym123",
			expectedCalls: func(m *MockGymHandler) {
				m.On("UpdateGym", mock.Anything, mock.Anything, "gym123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "activate gym",
			method: http.MethodPut,
			path:   "/gym123/activate",
			expectedCalls: func(m *MockGymHandler) {
				m.On("SetGymActive", mock.Anything, mock.Anything, "gym123", true).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "deactivate gym",
			method: http.MethodPut,
			path:   "/gym123/deactivate",
			expectedCalls: func(m *MockGymHandler) {
				m.On("SetGymActive", mock.Anything, mock.Anything, "gym123", false).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "delete gym",
			method: http.MethodDelete,
			path:   "/gym123",
			expectedCalls: func(m *MockGymHandler) {
				m.On("DeleteGym", mock.Anything, mock.Anything, "gym123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusNoContent)
				})
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedCalls(mockHandler)

			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			mockHandler.AssertExpectations(t)
		})
	}
}
