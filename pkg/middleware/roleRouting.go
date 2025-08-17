package middleware

import (
	"tesodev-korpes/pkg/customError"
	"tesodev-korpes/shared/config"

	"github.com/labstack/echo/v4"
)

func RoleRouting(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("userRole").(string)
			if !ok || userRole == "" {
				return customError.NewBadRequest(customError.EmptyRole)
			}

			internalPath, ok := cfg.RoleMapping[userRole]
			if !ok {
				return customError.NewForbidden(customError.ForbiddenAccess)
			}

			// Path’i config’teki internal path ile değiştir
			c.SetPath(internalPath)

			return next(c)
		}
	}
}
