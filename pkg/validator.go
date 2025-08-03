package pkg

import (
	"fmt"
	"net/http"
)

type ValidationErrorDetail struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type AppValidationError struct {
	*AppError
	Details []ValidationErrorDetail `json:"details"`
}

func ValidationFailed(details []ValidationErrorDetail, message string) *AppValidationError {
	baseErr := &AppError{
		HTTPStatus: http.StatusUnprocessableEntity, // 422, validasyon hataları için
		Code:       CodeValidation,
		Message:    message,
	}

	fmt.Println("error validatedesins")
	return &AppValidationError{
		AppError: baseErr,
		Details:  details,
	}
}
