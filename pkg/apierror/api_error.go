package apierror

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *APIError) Error() string { return e.Message }
func New(code, message string, err error) *APIError {
	return &APIError{Code: code, Message: message, Err: err}
}
