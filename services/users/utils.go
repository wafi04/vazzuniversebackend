package users

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrorCode adalah enum untuk kode kesalahan spesifik
type ErrorCode string

const (
	// Auth related errors
	ErrInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrUserNotFound       ErrorCode = "USER_NOT_FOUND"
	ErrInvalidPassword    ErrorCode = "INVALID_PASSWORD"
	ErrUsernameTaken      ErrorCode = "USERNAME_TAKEN"
	ErrEmailTaken         ErrorCode = "EMAIL_TAKEN"
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	
	// Payment related errors
	ErrInvalidPaymentCode ErrorCode = "INVALID_PAYMENT_CODE"
	ErrInsufficientFunds  ErrorCode = "INSUFFICIENT_FUNDS"
	ErrPaymentFailed      ErrorCode = "PAYMENT_FAILED"
	ErrInvalidAmount      ErrorCode = "INVALID_AMOUNT"
	
	// Validation errors
	ErrInvalidInput       ErrorCode = "INVALID_INPUT"
	ErrMissingField       ErrorCode = "MISSING_FIELD"
	ErrInvalidFormat      ErrorCode = "INVALID_FORMAT"
	
	// Generic errors
	ErrInternalServer     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrDatabaseError      ErrorCode = "DATABASE_ERROR"
)

// ResponseError adalah struktur standar untuk mengembalikan error API
type ResponseError struct {
	Timestamp   time.Time  `json:"timestamp"`
	Status      int        `json:"status"`      // HTTP status code
	Code        ErrorCode  `json:"code"`        // Application specific error code
	Message     string     `json:"message"`     // Human-readable error message
	Details     string     `json:"details,omitempty"` // Detailed error information if available
	Path        string     `json:"path,omitempty"`    // API endpoint path where error occurred
	RequestID   string     `json:"requestId,omitempty"` // For tracking purposes
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

// Beberapa fungsi pembuat error umum
func ErrUserNotFoundError(username string) *ResponseError {
	return NewResponseError(
		http.StatusNotFound,
		ErrUserNotFound,
		"User not found",
	).WithDetails("No user found with username: " + username)
}

func ErrInvalidCredentialsError() *ResponseError {
	return NewResponseError(
		http.StatusUnauthorized,
		ErrInvalidCredentials,
		"Invalid username or password",
	)
}

func ErrInvalidPaymentCodeError(paymentCode string) *ResponseError {
	return NewResponseError(
		http.StatusBadRequest,
		ErrInvalidPaymentCode,
		"Invalid payment code",
	).WithDetails("The payment code " + paymentCode + " is invalid or has expired")
}

func ErrInsufficientBalanceError() *ResponseError {
	return NewResponseError(
		http.StatusBadRequest,
		ErrInsufficientFunds,
		"Insufficient balance for this operation",
	)
}

func ErrUsernameTakenError(username string) *ResponseError {
	return NewResponseError(
		http.StatusConflict,
		ErrUsernameTaken,
		"Username already taken",
	).WithDetails("The username " + username + " is already in use")
}

func ErrEmailTakenError(email string) *ResponseError {
	return NewResponseError(
		http.StatusConflict,
		ErrEmailTaken,
		"Email already registered",
	).WithDetails("The email " + email + " is already registered")
}

func ErrInternalServerError(err error) *ResponseError {
	return NewResponseError(
		http.StatusInternalServerError,
		ErrInternalServer,
		"Internal server error occurred",
	).WithDetails(err.Error())
}

func HashPassword(password string) (string, error) {

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}
