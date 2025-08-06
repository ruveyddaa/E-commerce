package errorPackage

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
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

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, CodeNotFound, message)
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

func Internal(err error, message string) *AppError {
	log.Error("Internal error:", err)
	return Wrap(err, http.StatusInternalServerError, CodeInternalError, message)
}

func InvalidOrderStateWithStatus(action, currentStatus string) *AppError {
	message := fmt.Sprintf("Cannot %s order while it is in '%s' status", action, currentStatus)
	return New(http.StatusConflict, CodeOrderStateConflict, message)
}

func UnauthorizedInvalidLogin() *AppError {
	return New(http.StatusUnauthorized, CodeUnauthorized, "Invalid email or password")
}

func UnauthorizedInvalidToken() *AppError {
	return New(http.StatusUnauthorized, CodeUnauthorized, "Invalid or missing authorization token")
}
