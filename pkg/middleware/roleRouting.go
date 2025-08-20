package middleware

import (
	"errors"
	"fmt"
	"strings" // strings paketini import etmeyi unutma
	"tesodev-korpes/pkg/customError"
	"tesodev-korpes/shared/config"

	"github.com/labstack/echo/v4"
)

func RoleRouting(cfg config.Config, handlers map[string]echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("userMembership").(string)
			if !ok || userRole == "" {
				return customError.NewBadRequest(customError.EmptyRole)
			}

			internalPathTemplate, ok := cfg.RoleMapping[userRole]
			if !ok {
				return customError.NewForbidden(customError.ForbiddenAccess)
			}
			fmt.Println("Template Path:", internalPathTemplate)

			orderID := c.Param("id")
			internalPath := strings.Replace(internalPathTemplate, ":id", orderID, 1)

			fmt.Println("Final Path:", internalPath)

			targetHandler, found := handlers[internalPathTemplate]
			if !found {
				return customError.NewInternal(customError.OrderServiceError, errors.New("not exist endpoint"))
			}

			return targetHandler(c)
		}
	}
}
