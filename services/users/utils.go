package users

import (
	"net/http"

	"github.com/wafi04/vazzuniversebackend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Auth related errors
	ErrInvalidCredentials utils.ErrorCode = "INVALID_CREDENTIALS"
	ErrUserNotFound       utils.ErrorCode = "USER_NOT_FOUND"
	ErrInvalidPassword    utils.ErrorCode = "INVALID_PASSWORD"
	ErrUsernameTaken      utils.ErrorCode = "USERNAME_TAKEN"
	ErrEmailTaken         utils.ErrorCode = "EMAIL_TAKEN"
	ErrUnauthorized       utils.ErrorCode = "UNAUTHORIZED"
	
	// Payment related errors
	ErrInvalidPaymentCode utils.ErrorCode = "INVALID_PAYMENT_CODE"
	ErrInsufficientFunds  utils.ErrorCode = "INSUFFICIENT_FUNDS"
	ErrPaymentFailed      utils.ErrorCode = "PAYMENT_FAILED"
	ErrInvalidAmount      utils.ErrorCode = "INVALID_AMOUNT"
	
	// Validation errors
	ErrInvalidInput       utils.ErrorCode = "INVALID_INPUT"
	ErrMissingField       utils.ErrorCode = "MISSING_FIELD"
	ErrInvalidFormat      utils.ErrorCode = "INVALID_FORMAT"

	
	// Generic errors
	ErrInternalServer     utils.ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrDatabaseError      utils.ErrorCode = "DATABASE_ERROR"
)


func ErrUserNotFoundError(username string) *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusNotFound,
		ErrUserNotFound,
		"User not found",
	).WithDetails("No user found with username: " + username)
}


func ErrInvalidPasswordError() *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusUnauthorized,
		ErrInvalidPassword,
		"Invalid password",
	).WithDetails("The password provided is incorrect")
}

func ErrUnauthorizedError() *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusUnauthorized,
		ErrUnauthorized,
		"Unauthorized access",
	).WithDetails("You do not have permission to access this resource")
}

func ErrPaymentFailedError() *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusPaymentRequired,
		ErrPaymentFailed,
		"Payment failed",
	).WithDetails("The payment process encountered an error")
}

func ErrInvalidAmountError() *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusBadRequest,
		ErrInvalidAmount,
		"Invalid amount",
	).WithDetails("The amount provided is invalid or exceeds the limit")
}

func ErrInvalidCredentialsError() *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusUnauthorized,
		ErrInvalidCredentials,
		"Invalid username or password",
	)
}

func ErrInvalidPaymentCodeError(paymentCode string) *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusBadRequest,
		ErrInvalidPaymentCode,
		"Invalid payment code",
	).WithDetails("The payment code " + paymentCode + " is invalid or has expired")
}

func ErrInsufficientBalanceError() *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusBadRequest,
		ErrInsufficientFunds,
		"Insufficient balance for this operation",
	)
}

func ErrUsernameTakenError(username string) *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusConflict,
		ErrUsernameTaken,
		"Username already taken",
	).WithDetails("The username " + username + " is already in use")
}

func ErrEmailTakenError(email string) *utils.ResponseError {
	return utils.NewResponseError(
		http.StatusConflict,
		ErrEmailTaken,
		"Email already registered",
	).WithDetails("The email " + email + " is already registered")
}

func ErrInternalServerError(err error) *utils.ResponseError {
	return utils.NewResponseError(
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
