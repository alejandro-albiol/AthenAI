package response

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

// APIResponse is a generic API response wrapper.
type APIResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// WriteAPISuccess writes a standardized JSON success response.
func WriteAPISuccess(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse[any]{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// WriteAPICreated writes a standardized JSON success response with 201 status code.
func WriteAPICreated(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(APIResponse[any]{
		Status:  "success",
		Message: message,
		Data:    map[string]any{"id": data},
	})
}

func WriteAPIError(w http.ResponseWriter, apiErr *apierror.APIError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	status := http.StatusBadRequest
	switch apiErr.Code {
	case errorcode_enum.CodeConflict:
		status = http.StatusConflict
	case errorcode_enum.CodeInternal:
		status = http.StatusInternalServerError
	case errorcode_enum.CodeNotFound:
		status = http.StatusNotFound
	case errorcode_enum.CodeUnauthorized:
		status = http.StatusUnauthorized
	case errorcode_enum.CodeForbidden:
		status = http.StatusForbidden
	}
	w.WriteHeader(status)

	// Create error data based on environment
	data := map[string]any{"code": apiErr.Code}

	// Only include detailed error information in development mode
	if isDevelopmentMode() && apiErr.Err != nil {
		data["error"] = apiErr.Err.Error()
	}
	// In production, only return the code - no error details

	json.NewEncoder(w).Encode(APIResponse[any]{
		Status:  "error",
		Message: apiErr.Message,
		Data:    data,
	})
}

// isDevelopmentMode checks if the application is running in development mode
func isDevelopmentMode() bool {
	env := os.Getenv("APP_ENV")
	return env == "development" || env == "dev"
}
