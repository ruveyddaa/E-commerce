package pkg

import (
	"fmt"
	"net/http"
	"tesodev-korpes/pkg/customError"
)

type ValidationErrorDetail struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type AppValidationError struct {
	*customError.AppError
	Details []ValidationErrorDetail `json:"details"`
}

func ValidationFailed(details []ValidationErrorDetail, message string) *AppValidationError {
	baseErr := &customError.AppError{
		HTTPStatus: http.StatusUnprocessableEntity, // 422, for validation errors
		Message:    message,
	}

	fmt.Println("in here error validation")
	return &AppValidationError{
		AppError: baseErr,
		Details:  details,
	}
}
