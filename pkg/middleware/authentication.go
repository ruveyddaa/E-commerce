package middleware

import (
	"strings"
	"tesodev-korpes/pkg/auth"
	"tesodev-korpes/pkg/customError"

	"github.com/labstack/echo/v4"
)

type SkipperFunc func(c echo.Context) bool

func Authentication(skipper SkipperFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper != nil && skipper(c) {
				return next(c)
			}
			const bearerPrefix = "Bearer "

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
				return customError.NewUnauthorized(customError.MissingAuthToken)
			}

			tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

			claims, err := auth.VerifyJWT(tokenStr)
			if err != nil {
				return customError.NewUnauthorized(customError.MissingAuthToken)
			}

			c.Set("userID", claims.ID)
			return next(c)
		}
	}
}
