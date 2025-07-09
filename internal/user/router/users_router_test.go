package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserHandler struct {
	mock.Mock
}

func (m *MockUserHandler) RegisterUser(w http.ResponseWriter, r *http.Request, gymID string) {
	m.Called(w, r, gymID)
}

func (m *MockUserHandler) GetUserByID(w http.ResponseWriter, gymID, id string) {
	m.Called(w, gymID, id)
}

func (m *MockUserHandler) GetUserByUsername(w http.ResponseWriter, gymID, username string) {
	m.Called(w, gymID, username)
}

func (m *MockUserHandler) GetUserByEmail(w http.ResponseWriter, gymID, email string) {
	m.Called(w, gymID, email)
}

func (m *MockUserHandler) GetAllUsers(w http.ResponseWriter, gymID string) {
	m.Called(w, gymID)
}

func (m *MockUserHandler) UpdateUser(w http.ResponseWriter, gymID, id string, user dto.UserUpdateDTO) {
	m.Called(w, gymID, id, user)
}

func (m *MockUserHandler) DeleteUser(w http.ResponseWriter, gymID, id string) {
	m.Called(w, gymID, id)
}

func (m *MockUserHandler) VerifyUser(w http.ResponseWriter, gymID, id string) {
	m.Called(w, gymID, id)
}

func (m *MockUserHandler) SetUserActive(w http.ResponseWriter, gymID, id string, active bool) {
	m.Called(w, gymID, id, active)
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
			name:   "missing gym ID header",
			method: http.MethodGet,
			path:   "/",
			gymID:  "",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				// No mock setup needed, middleware should reject
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "get all users",
			method: http.MethodGet,
			path:   "/",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("GetAllUsers", w, "gym123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "get user by id",
			method: http.MethodGet,
			path:   "/user123",
			gymID:  "gym123",
			setupMock: func(m *MockUserHandler, w *httptest.ResponseRecorder) {
				m.On("GetUserByID", w, "gym123", "user123").Run(func(args mock.Arguments) {
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
				m.On("RegisterUser", w, mock.Anything, "gym123").Run(func(args mock.Arguments) {
					w := args.Get(0).(http.ResponseWriter)
					w.WriteHeader(http.StatusCreated)
				}).Once()
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := new(MockUserHandler)
			router := NewUsersRouter(mockHandler)

			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
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
