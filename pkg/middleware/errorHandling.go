package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"tesodev-korpes/pkg"
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
		case http.StatusBadRequest:
			return pkg.BadRequest(fmt.Sprintf("%v", httpErr.Message))
		case http.StatusInternalServerError:
			return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceServiceCode500301])
		default:
			return pkg.Wrap(httpErr, httpErr.Code, pkg.CodeInternalFrameworkError, pkg.InternalServerErrorMessages[pkg.ResourceFrameworkCode500401])
		}
	}

	return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceServiceCode500301])
}

func buildAPIResponse(err *pkg.AppError, c echo.Context) APIErrorResponse {
	corralationID := c.Response().Header().Get(echo.HeaderXCorrelationID)
	if corralationID == "" {
		corralationID = "not-available"
	}

	var resp APIErrorResponse
	resp.Error.Code = err.Code
	resp.Error.Message = err.Message
	resp.CorralationID = corralationID
	resp.Timestamp = time.Now().UTC().Format(time.RFC3339)
	return resp
}
