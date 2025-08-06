package middleware

import (
	"fmt"
	"net/http"
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
		customer, err := fetchCustomerByID("http://localhost:8001", claims.ID)
		if err != nil || customer == nil {
			fmt.Println(err)
			fmt.Println(customer)
			return echo.NewHTTPError(http.StatusUnauthorized, "Müşteri bilgisi alınamadı")
		}

		c.Set("userID", claims.ID)
		c.Set("userRole", customer.Role)

		return next(c)
	}
}
