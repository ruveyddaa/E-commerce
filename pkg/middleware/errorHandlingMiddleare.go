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
			// A defer/recover block to catch panics.
			// This prevents the server from crashing in case of an unexpected panic
			// in any handler and returns a standard error message.
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("panic: %v", r)
					}

					// Convert the panic into a standard internal server error
					appErr := pkg.Internal(err)

					// Send the JSON response to the client
					if !c.Response().Committed {
						_ = c.JSON(appErr.HTTPStatus, buildAPIResponse(appErr, c))
					}
				}
			}()

			// Run the next middleware or handler.
			err := next(c)
			if err == nil {
				return nil
			}

			// Convert the incoming error to the AppError format.
			appErr := toAppError(err)

			// If the response has already been sent (e.g., for a WebSocket), do not send it again.
			if c.Response().Committed {
				return nil
			}

			// Create the standard API response and send it to the client.
			return c.JSON(appErr.HTTPStatus, buildAPIResponse(appErr, c))
		}
	}
}

// toAppError converts any `error` type into our standard `AppError` type.
func toAppError(err error) *pkg.AppError {
	var appErr *pkg.AppError

	// If the error is already an AppError, return it directly.
	if errors.As(err, &appErr) {
		return appErr
	}

	// If the error is of Echo's own type, we convert it to the AppError type.
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		switch httpErr.Code {
		case http.StatusNotFound:
			return pkg.NotFound()
		case http.StatusBadRequest:
			return pkg.BadRequest(fmt.Sprintf("%v", httpErr.Message))
		default:
			// Wrap other Echo errors as a generic internal error.
			return pkg.Wrap(httpErr, httpErr.Code, "INTERNAL_FRAMEWORK_ERROR", "A framework-related error occurred.")
		}
	}

	// For any other type of error, wrap it as a generic internal error.
	return pkg.Internal(err)
}

// buildAPIResponse constructs the final JSON response body from an AppError.
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