package pkg

import (
	"fmt"
	"net/http"
)

const (
	CodeUnknown            = "UNKNOWN"
	CodeInvalidInput       = "INVALID_INPUT"
	CodeValidation         = "VALIDATION_FAILED"
	CodeNotFound           = "NOT_FOUND"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeConflict           = "CONFLICT"
	CodeInternalError      = "INTERNAL_ERROR"
	CodeServiceDown        = "SERVICE_UNAVAILABLE"
	CodeOrderStateConflict = "ORDER_STATE_CONFLICT"
)

type AppError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %s, message: %s, underlying_error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(status int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
	}
}

func Wrap(err error, status int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func NotFound() *AppError {
	return New(http.StatusNotFound, CodeNotFound, "The requested resource was not found.")
}

func ValidationFailed(message string) *AppError {
	if message == "" {
		message = "Input validation failed."
	}
	return New(http.StatusUnprocessableEntity, CodeValidation, message)
}

func BadRequest(message string) *AppError {
	if message == "" {
		message = "Invalid or missing parameter."
	}
	return New(http.StatusBadRequest, CodeInvalidInput, message)
}

func Unauthorized() *AppError {
	return New(http.StatusUnauthorized, CodeUnauthorized, "You are not authorized to perform this action.")
}

func Forbidden() *AppError {
	return New(http.StatusForbidden, CodeForbidden, "You do not have permission to access this resource.")
}

func Internal(err error) *AppError {
	return Wrap(err, http.StatusInternalServerError, CodeInternalError, "An unexpected error occurred in the system.")
}

func InvalidOrderStateWithStatus(action, currentStatus string) *AppError {
	message := fmt.Sprintf("Cannot %s order while it is in '%s' status", action, currentStatus)
	return New(http.StatusConflict, CodeOrderStateConflict, message)
}
