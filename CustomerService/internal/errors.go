package internal

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func NewInternal(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func NewValidation(msg string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: msg}
}

func Respond(c echo.Context, err error, fallbackMsg string) error {
	if appErr, ok := err.(*AppError); ok {
		return c.JSON(appErr.Code, echo.Map{"error": appErr.Message})
	}
	return c.JSON(http.StatusInternalServerError, echo.Map{"error": fallbackMsg})
}
