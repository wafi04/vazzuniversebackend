package response

import "time"

type ErrorCode string

type ResponseError struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Path      string    `json:"path,omitempty"`
	RequestID string    `json:"requestId,omitempty"`
}

// NewResponseError membuat instance ResponseError baru
func NewResponseError(status int, code ErrorCode, message string) *ResponseError {
	return &ResponseError{
		Timestamp: time.Now(),
		Status:    status,
		Code:      code,
		Message:   message,
	}
}

// WithDetails menambahkan detail ke ResponseError
func (e *ResponseError) WithDetails(details string) *ResponseError {
	e.Details = details
	return e
}

// WithPath menambahkan path ke ResponseError
func (e *ResponseError) WithPath(path string) *ResponseError {
	e.Path = path
	return e
}

// WithRequestID menambahkan requestID ke ResponseError
func (e *ResponseError) WithRequestID(requestID string) *ResponseError {
	e.RequestID = requestID
	return e
}
