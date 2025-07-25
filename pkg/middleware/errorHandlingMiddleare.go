package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"tesodev-korpes/pkg"
	"time"

	"github.com/labstack/echo/v4"
)

// APIErrorResponse is the standard JSON structure to be returned to the client.
type APIErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

// ErrorHandler creates the main error handler middleware for the project using Echo.
func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			appErr := toAppError(err)

			if c.Response().Committed {
				return nil
			}

			return c.JSON(appErr.HTTPStatus, buildAPIResponse(appErr, c))
		}
	}
}

func toAppError(err error) *pkg.AppError {
	var appErr *pkg.AppError

	if errors.As(err, &appErr) {
		return appErr
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		switch httpErr.Code {
		case http.StatusNotFound:
			return pkg.NotFound()
		case http.StatusBadRequest:
			return pkg.BadRequest(fmt.Sprintf("%v", httpErr.Message))
		default:
			return pkg.Wrap(httpErr, httpErr.Code, "INTERNAL_FRAMEWORK_ERROR", "A framework-related error occurred.")
		}
	}

	return pkg.Internal(err)
}

func buildAPIResponse(err *pkg.AppError, c echo.Context) APIErrorResponse {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = "not-available"
	}

	var resp APIErrorResponse
	resp.Error.Code = err.Code
	resp.Error.Message = err.Message
	resp.RequestID = requestID
	resp.Timestamp = time.Now().UTC().Format(time.RFC3339)
	return resp
}
