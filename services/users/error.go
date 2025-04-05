package users

import (
	"net/http"

	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Auth related errors
	ErrInvalidCredentials response.ErrorCode = "INVALID_CREDENTIALS"
	ErrUserNotFound       response.ErrorCode = "USER_NOT_FOUND"
	ErrInvalidPassword    response.ErrorCode = "INVALID_PASSWORD"
	ErrUsernameTaken      response.ErrorCode = "USERNAME_TAKEN"
	ErrEmailTaken         response.ErrorCode = "EMAIL_TAKEN"
	ErrUnauthorized       response.ErrorCode = "UNAUTHORIZED"

	// Payment related errors
	ErrInvalidPaymentCode response.ErrorCode = "INVALID_PAYMENT_CODE"
	ErrInsufficientFunds  response.ErrorCode = "INSUFFICIENT_FUNDS"
	ErrPaymentFailed      response.ErrorCode = "PAYMENT_FAILED"
	ErrInvalidAmount      response.ErrorCode = "INVALID_AMOUNT"

	// Validation errors
	ErrInvalidInput  response.ErrorCode = "INVALID_INPUT"
	ErrMissingField  response.ErrorCode = "MISSING_FIELD"
	ErrInvalidFormat response.ErrorCode = "INVALID_FORMAT"

	// Generic errors
	ErrInternalServer response.ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrDatabaseError  response.ErrorCode = "DATABASE_ERROR"
)

func ErrUserNotFoundError(username string) *response.ResponseError {
	return response.NewResponseError(
		http.StatusNotFound,
		ErrUserNotFound,
		"User not found",
	).WithDetails("No user found with username: " + username)
}

func ErrInvalidPasswordError() *response.ResponseError {
	return response.NewResponseError(
		http.StatusUnauthorized,
		ErrInvalidPassword,
		"Invalid password",
	).WithDetails("The password provided is incorrect")
}

func ErrUnauthorizedError() *response.ResponseError {
	return response.NewResponseError(
		http.StatusUnauthorized,
		ErrUnauthorized,
		"Unauthorized access",
	).WithDetails("You do not have permission to access this resource")
}

func ErrPaymentFailedError() *response.ResponseError {
	return response.NewResponseError(
		http.StatusPaymentRequired,
		ErrPaymentFailed,
		"Payment failed",
	).WithDetails("The payment process encountered an error")
}

func ErrInvalidAmountError() *response.ResponseError {
	return response.NewResponseError(
		http.StatusBadRequest,
		ErrInvalidAmount,
		"Invalid amount",
	).WithDetails("The amount provided is invalid or exceeds the limit")
}

func ErrInvalidCredentialsError() *response.ResponseError {
	return response.NewResponseError(
		http.StatusUnauthorized,
		ErrInvalidCredentials,
		"Invalid username or password",
	)
}

func ErrInvalidPaymentCodeError(paymentCode string) *response.ResponseError {
	return response.NewResponseError(
		http.StatusBadRequest,
		ErrInvalidPaymentCode,
		"Invalid payment code",
	).WithDetails("The payment code " + paymentCode + " is invalid or has expired")
}

func ErrInsufficientBalanceError() *response.ResponseError {
	return response.NewResponseError(
		http.StatusBadRequest,
		ErrInsufficientFunds,
		"Insufficient balance for this operation",
	)
}

func ErrUsernameTakenError(username string) *response.ResponseError {
	return response.NewResponseError(
		http.StatusConflict,
		ErrUsernameTaken,
		"Username already taken",
	).WithDetails("The username " + username + " is already in use")
}

func ErrEmailTakenError(email string) *response.ResponseError {
	return response.NewResponseError(
		http.StatusConflict,
		ErrEmailTaken,
		"Email already registered",
	).WithDetails("The email " + email + " is already registered")
}

func ErrInternalServerError(err error) *response.ResponseError {
	return response.NewResponseError(
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

func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
