package middleware

import (
	"fmt"
	"strings" // strings paketini import etmeyi unutma
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

			internalPathTemplate, ok := cfg.RoleMapping[userRole]
			if !ok {
				return customError.NewForbidden(customError.ForbiddenAccess)
			}
			fmt.Println("Template Path:", internalPathTemplate) // Örn: /internal/price/non-premium/:id

			// ÖNEMLİ DEĞİŞİKLİK: ID'yi alıp yeni yola ekliyoruz.
			orderID := c.Param("id")
			internalPath := strings.Replace(internalPathTemplate, ":id", orderID, 1)

			fmt.Println("Final Path:", internalPath) // Örn: /internal/price/non-premium/e18743ec...

			c.SetPath(internalPath)

			return next(c)
		}
	}
}
