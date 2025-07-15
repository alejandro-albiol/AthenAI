package response_test

import (
	"errors"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestWriteAPIError_SecurityModes(t *testing.T) {
	// Test error with sensitive database information
	sensitiveErr := apierror.New(
		errorcode_enum.CodeUnauthorized,
		"Invalid refresh token",
		// Simulate a database error that reveals internal structure
		errors.New("sql: no rows in result set"),
	)

	t.Run("development mode shows detailed errors", func(t *testing.T) {
		// Set development environment
		os.Setenv("APP_ENV", "dev")
		defer os.Unsetenv("APP_ENV")

		w := httptest.NewRecorder()
		response.WriteAPIError(w, sensitiveErr)

		body := w.Body.String()
		assert.Contains(t, body, `"code":"UNAUTHORIZED"`)
		assert.Contains(t, body, `"error"`)                    // Should contain error details
		assert.Contains(t, body, "sql: no rows in result set") // Should show database details
	})

	t.Run("production mode hides sensitive errors", func(t *testing.T) {
		// Set production environment
		os.Setenv("APP_ENV", "prod")
		defer os.Unsetenv("APP_ENV")

		w := httptest.NewRecorder()
		response.WriteAPIError(w, sensitiveErr)

		body := w.Body.String()
		assert.Contains(t, body, `"code":"UNAUTHORIZED"`)
		assert.NotContains(t, body, `"error":`)                   // Should NOT contain error field
		assert.NotContains(t, body, "sql: no rows in result set") // Should NOT show database details

		// Verify only code is present in data
		assert.Contains(t, body, `"data":{"code":"UNAUTHORIZED"}`)
	})

	t.Run("unset environment defaults to production", func(t *testing.T) {
		// Ensure no APP_ENV is set (defaults to production)
		os.Unsetenv("APP_ENV")

		w := httptest.NewRecorder()
		response.WriteAPIError(w, sensitiveErr)

		body := w.Body.String()
		assert.Contains(t, body, `"code":"UNAUTHORIZED"`)
		assert.NotContains(t, body, `"error":`)                   // Should NOT contain error field (production mode)
		assert.NotContains(t, body, "sql: no rows in result set") // Should NOT show database details
	})
}

func TestWriteAPIError_StatusCodes(t *testing.T) {
	testCases := []struct {
		code           string
		expectedStatus int
	}{
		{errorcode_enum.CodeBadRequest, 400},
		{errorcode_enum.CodeUnauthorized, 401},
		{errorcode_enum.CodeForbidden, 403},
		{errorcode_enum.CodeNotFound, 404},
		{errorcode_enum.CodeConflict, 409},
		{errorcode_enum.CodeInternal, 500},
	}

	for _, tc := range testCases {
		t.Run(tc.code, func(t *testing.T) {
			apiErr := apierror.New(tc.code, "Test message", nil)
			w := httptest.NewRecorder()

			response.WriteAPIError(w, apiErr)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.code)
		})
	}
}
