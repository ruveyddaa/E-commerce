package middleware

import (
	"errors"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/errorPackage"
	"time"

	"github.com/labstack/echo/v4"
)

type APIErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	CorralationID string `json:"corralation_id"`
	Timestamp     string `json:"timestamp"`
}

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var validationErr *pkg.AppValidationError
			if errors.As(err, &validationErr) {
				return c.JSON(validationErr.HTTPStatus, validationErr)
			}

			appErr := toAppError(err)

			if c.Response().Committed {
				return nil
			}

			return c.JSON(appErr.HTTPStatus, buildAPIResponse(appErr, c))
		}
	}
}

func toAppError(err error) *errorPackage.AppError {
	var appErr *errorPackage.AppError

	if errors.As(err, &appErr) {
		return appErr
	}

	return errorPackage.NewInternal("500001", err)
}

func buildAPIResponse(err *errorPackage.AppError, c echo.Context) APIErrorResponse {
	correlationID := c.Response().Header().Get(echo.HeaderXCorrelationID)
	if correlationID == "" {
		correlationID = "not-available"
	}

	var resp APIErrorResponse
	resp.Error.Code = err.Code
	resp.Error.Message = err.Message
	resp.CorralationID = correlationID
	resp.Timestamp = time.Now().UTC().Format(time.RFC3339)
	return resp
}
