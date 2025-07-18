package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
	"github.com/alejandro-albiol/athenai/internal/user/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserHandler struct {
	mock.Mock
}

func (m *MockUserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockUserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	m.Called(w, r, id)
}

func (m *MockUserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request, username string) {
	m.Called(w, r, username)
}

func (m *MockUserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request, email string) {
	m.Called(w, r, email)
}

func (m *MockUserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockUserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id string, user dto.UserUpdateDTO) {
	m.Called(w, r, id, user)
}

func (m *MockUserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	m.Called(w, r, id)
}

func (m *MockUserHandler) VerifyUser(w http.ResponseWriter, r *http.Request, id string) {
	m.Called(w, r, id)
}

func (m *MockUserHandler) SetUserActive(w http.ResponseWriter, r *http.Request, id string, active bool) {
	m.Called(w, r, id, active)
}

func TestUserRoutes(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		gymID          string
		body           interface{}
		setupMock      func(*MockUserHandler, *httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name:   "get all users",
			method: http.MethodGet,
			path:   "/",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("GetAllUsers", w, mock.Anything).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "get all users with empty gym ID",
			method: http.MethodGet,
			path:   "/",
			gymID:  "",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				// Handler will be called but should handle empty/invalid gym context gracefully
				m.On("GetAllUsers", w, mock.Anything).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusBadRequest)
				}).Maybe()
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "get user by id",
			method: http.MethodGet,
			path:   "/user123",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("GetUserByID", w, mock.Anything, "user123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "get user by username",
			method: http.MethodGet,
			path:   "/username/testuser",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("GetUserByUsername", w, mock.Anything, "testuser").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "get user by email",
			method: http.MethodGet,
			path:   "/email/test@test.com",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("GetUserByEmail", w, mock.Anything, "test@test.com").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "register user",
			method: http.MethodPost,
			path:   "/",
			gymID:  "gym123",
			body: dto.UserCreationDTO{
				Username: "testuser",
				Email:    "test@test.com",
				Password: "password123",
				Role:     userrole_enum.User,
			},
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("RegisterUser", w, mock.Anything).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusCreated)
				}).Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "register user with invalid JSON",
			method: http.MethodPost,
			path:   "/",
			gymID:  "gym123",
			body:   `{"invalid": json}`, // Invalid JSON will be handled by handler
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("RegisterUser", w, mock.Anything).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusBadRequest)
				}).Once()
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "update user",
			method: http.MethodPut,
			path:   "/user123",
			gymID:  "gym123",
			body: dto.UserUpdateDTO{
				Username: "updateduser",
				Email:    "updated@test.com",
			},
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("UpdateUser", w, mock.Anything, "user123", mock.AnythingOfType("dto.UserUpdateDTO")).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "update user with invalid JSON",
			method: http.MethodPut,
			path:   "/user123",
			gymID:  "gym123",
			body:   `{"invalid": json}`, // Invalid JSON
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				// No mock setup needed as router should handle JSON error before calling handler
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "delete user",
			method: http.MethodDelete,
			path:   "/user123",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("DeleteUser", w, mock.Anything, "user123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusNoContent)
				})
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "verify user",
			method: http.MethodPost,
			path:   "/user123/verify",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("VerifyUser", w, mock.Anything, "user123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "set user active",
			method: http.MethodPost,
			path:   "/user123/active",
			gymID:  "gym123",
			body:   map[string]bool{"active": true},
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("SetUserActive", w, mock.Anything, "user123", true).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "set user inactive",
			method: http.MethodPost,
			path:   "/user123/active",
			gymID:  "gym123",
			body:   map[string]bool{"active": false},
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("SetUserActive", w, mock.Anything, "user123", false).Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "set user active with invalid JSON",
			method: http.MethodPost,
			path:   "/user123/active",
			gymID:  "gym123",
			body:   `{"invalid": json}`, // Invalid JSON
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				// No mock setup needed as router should handle JSON error before calling handler
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := new(MockUserHandler)
			router := router.NewUsersRouter(mockHandler)

			var body []byte
			if tt.body != nil {
				switch v := tt.body.(type) {
				case string:
					body = []byte(v)
				default:
					body, _ = json.Marshal(tt.body)
				}
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
			if tt.gymID != "" {
				req.Header.Set("X-Gym-ID", tt.gymID)
			}

			w := httptest.NewRecorder()
			tt.setupMock(mockHandler, w)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockHandler.AssertExpectations(t)
		})
	}
}
