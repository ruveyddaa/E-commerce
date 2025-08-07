package middleware

import (
	"strings"
	"tesodev-korpes/pkg/auth"
	"tesodev-korpes/pkg/errorPackage"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		const bearerPrefix = "Bearer "

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
			return errorPackage.UnauthorizedInvalidToken()
		}

		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := auth.VerifyJWT(tokenStr)
		if err != nil {
			return errorPackage.UnauthorizedInvalidToken()
		}

		c.Set("userID", claims.ID)

		return next(c)
	}
}
