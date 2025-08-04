package pkg

import (
	"fmt"
	"net/http"
	"tesodev-korpes/pkg/errorPackage"
)

type ValidationErrorDetail struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type AppValidationError struct {
	*errorPackage.AppError
	Details []ValidationErrorDetail `json:"details"`
}

func ValidationFailed(details []ValidationErrorDetail, message string) *AppValidationError {
	baseErr := &errorPackage.AppError{
		HTTPStatus: http.StatusUnprocessableEntity, // 422, for validation errors
		Code:       errorPackage.CodeValidation,
		Message:    message,
	}

	fmt.Println("in here error validation")
	return &AppValidationError{
		AppError: baseErr,
		Details:  details,
	}
}
