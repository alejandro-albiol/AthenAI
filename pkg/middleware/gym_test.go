package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireGymID(t *testing.T) {
	testCases := []struct {
		name         string
		gymID        string
		expectStatus int
		expectCalled bool
	}{
		{
			name:         "valid gym ID",
			gymID:        "gym123",
			expectStatus: http.StatusOK,
			expectCalled: true,
		},
		{
			name:         "missing gym ID",
			gymID:        "",
			expectStatus: http.StatusBadRequest,
			expectCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Track if next handler was called
			nextCalled := false
			var receivedGymID string

			// Mock next handler
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				receivedGymID = GetGymID(r)
				w.WriteHeader(http.StatusOK)
			})

			// Create middleware
			middleware := RequireGymID(nextHandler)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.gymID != "" {
				req.Header.Set("X-Gym-ID", tc.gymID)
			}

			w := httptest.NewRecorder()

			// Execute middleware
			middleware.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectStatus, w.Code)
			assert.Equal(t, tc.expectCalled, nextCalled)

			if tc.expectCalled {
				assert.Equal(t, tc.gymID, receivedGymID)
			}
		})
	}
}

func TestGetGymID(t *testing.T) {
	testCases := []struct {
		name     string
		setupReq func() *http.Request
		expected string
	}{
		{
			name: "get from context",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				ctx := req.Context()
				ctx = context.WithValue(ctx, GymIDKey, "gym-from-context")
				return req.WithContext(ctx)
			},
			expected: "gym-from-context",
		},
		{
			name: "fallback to header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.Header.Set("X-Gym-ID", "gym-from-header")
				return req
			},
			expected: "gym-from-header",
		},
		{
			name: "no gym ID",
			setupReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/test", nil)
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := tc.setupReq()
			result := GetGymID(req)
			assert.Equal(t, tc.expected, result)
		})
	}
}
