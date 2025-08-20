package middleware

import (
	"tesodev-korpes/pkg/customError"
	"tesodev-korpes/shared/config"

	"github.com/labstack/echo/v4"
)

// func AuthorizationMiddleware(allowedRoles []string) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			userRole, ok := c.Get("userRole").(string)
// 			fmt.Println(userRole)
// 			if !ok || userRole == "" {
// 				return customError.NewBadRequest(customError.EmptyRole)
// 			}

// 			for _, allowed := range allowedRoles {
// 				if strings.EqualFold(userRole, allowed) {
// 					return next(c)
// 				}
// 			}

// 			return customError.NewForbidden(customError.ForbiddenAccess)
// 		}
// 	}
// }

func AuthorizationMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()

			allowedRoles, found := cfg.EndpointRoles[path]
			if !found {
				return next(c)
			}
			userRole := c.Get("userRole")

			for _, role := range allowedRoles {
				if role == userRole {
					return next(c)
				}
			}

			return customError.NewForbidden(customError.ForbiddenAccess)
		}
	}
}
