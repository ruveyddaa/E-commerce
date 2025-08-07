package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthorizationMiddleware(allowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get("userID").(string)
			if !ok || userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Kullanıcı doğrulanamadı")
			}
			userRole, ok := c.Get("userRole").(string)
			fmt.Println(userRole)
			if !ok || userRole == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Role alınımadı")
			}

			for _, allowed := range allowedRoles {
				if strings.EqualFold(userRole, allowed) {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Bu işlemi yapmaya yetkiniz yok")
		}
	}
}
