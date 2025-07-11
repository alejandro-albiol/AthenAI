package apierror

import (
	"encoding/json"
	"net/http"

	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *APIError) Error() string { return e.Message }
func New(code, message string, err error) *APIError {
	return &APIError{Code: code, Message: message, Err: err}
}

// WriteAPIError writes a standardized API error response with the correct HTTP status code.
// It maps APIError codes to HTTP status codes and formats the response consistently.
func WriteAPIError(w http.ResponseWriter, apiErr *APIError) {
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
	data := map[string]any{"code": apiErr.Code}
	if apiErr.Err != nil {
		data["error"] = apiErr.Err.Error()
	}
	json.NewEncoder(w).Encode(response.APIResponse[any]{
		Status:  "error",
		Message: apiErr.Message,
		Data:    data,
	})
}
