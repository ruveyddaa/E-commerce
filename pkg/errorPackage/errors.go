package errorPackage

import (
	"fmt"

	"github.com/labstack/gommon/log"
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

func newAppError(typeCode string, err error, args ...interface{}) *AppError {
	details, ok := ErrorTypeCode[typeCode]
	if !ok {
		log.Errorf("FATAL: Undefined error code requested: %s", typeCode)
		details = ErrorTypeCode["500001"] // Tanımsız kod istenirse genel bir iç hata dön
		typeCode = "500001"
	}

	msg := details.Message
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return &AppError{
		HTTPStatus: details.StatusCode,
		Code:       typeCode,
		Message:    msg,
		Err:        err,
	}
}

func NewNotFound(typeCode string) *AppError {
	return newAppError(typeCode, nil)
}

func NewBadRequest(typeCode string) *AppError {
	return newAppError(typeCode, nil)
}

func NewUnauthorized(typeCode string) *AppError {
	return newAppError(typeCode, nil)
}

func NewForbidden(typeCode string) *AppError {
	return newAppError(typeCode, nil)
}

func NewConflict(typeCode string, args ...interface{}) *AppError {
	return newAppError(typeCode, nil, args...)
}

func NewUnprocessableEntity(typeCode string, err error) *AppError {
	return newAppError(typeCode, err)
}

func NewInternal(typeCode string, err error) *AppError {
	log.Errorf("Internal error occurred with code %s: %v", typeCode, err)
	return newAppError(typeCode, err)
}

func NewValidation(typeCode string) *AppError {
	return newAppError(typeCode, nil)
}

