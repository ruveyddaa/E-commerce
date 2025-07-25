package pkg

import (
	"fmt"
	"net/http"
)

// Error Codes - Machine-readable, consistent error identifiers.
// These codes allow the client (frontend) to take specific actions based on the incoming error.
const (
	CodeUnknown         = "UNKNOWN"
	CodeInvalidInput    = "INVALID_INPUT"
	CodeValidation      = "VALIDATION_FAILED"
	CodeNotFound        = "NOT_FOUND"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeConflict        = "CONFLICT" // e.g., trying to create a user that already exists
	CodeInternalError   = "INTERNAL_ERROR"
	CodeServiceDown     = "SERVICE_UNAVAILABLE"
)

// AppError is the struct that holds all the rich error information within the application.
// This struct is used for logging and debugging.
type AppError struct {
	// HTTPStatus is the HTTP status code to be sent to the client.
	HTTPStatus int `json:"-"`

	// Code is the machine-readable error code (e.g., "NOT_FOUND").
	Code string `json:"code"`

	// Message is the human-readable error message that can be safely displayed to the client.
	Message string `json:"message"`

	// Err is the original error we are wrapping. It is never shown to the client.
	// It is only used for logging and internal error analysis.
	Err error `json:"-"`
}

// Error implements the error interface. It provides a rich output when called during logging.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %s, message: %s, underlying_error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

// Unwrap provides compatibility with Go's standard `errors.Is` and `errors.As` functions.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError. Used when there is no original error.
func New(status int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
	}
}

// Wrap wraps an existing error with an AppError. This is the most frequently used function to avoid losing the original error.
func Wrap(err error, status int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// NotFound creates a new error for situations where a resource is not found.
func NotFound() *AppError {
	return New(http.StatusNotFound, CodeNotFound, "The requested resource was not found.")
}

// ValidationFailed creates a new error for input validation failures.
func ValidationFailed(message string) *AppError {
	if message == "" {
		message = "Input validation failed."
	}
	return New(http.StatusUnprocessableEntity, CodeValidation, message)
}

// BadRequest creates a new error for invalid requests from the client.
func BadRequest(message string) *AppError {
	if message == "" {
		message = "Invalid or missing parameter."
	}
	return New(http.StatusBadRequest, CodeInvalidInput, message)
}

// Unauthorized creates a new error for unauthorized access attempts.
func Unauthorized() *AppError {
	return New(http.StatusUnauthorized, CodeUnauthorized, "You are not authorized to perform this action.")
}

// Forbidden creates a new error for attempts to access a forbidden resource.
func Forbidden() *AppError {
	return New(http.StatusForbidden, CodeForbidden, "You do not have permission to access this resource.")
}

// Internal creates a new error for unexpected internal system errors.
func Internal(err error) *AppError {
	return Wrap(err, http.StatusInternalServerError, CodeInternalError, "An unexpected error occurred in the system.")
}
