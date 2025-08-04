package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"tesodev-korpes/CustomerService/authentication"
)

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		skipConditions := []struct {
			Method string
			Path   string
		}{
			{Method: "POST", Path: "/login"},
			{Method: "POST", Path: "/customer"},
			{Method: "GET", Path: "/verify"},
			{Method: "GET", Path: "/swagger/*"},
			{Method: "GET", Path: "/customers"},
			{Method: "GET", Path: "/docs/*"},
			{Method: "GET", Path: "/order/swagger/*"},
			{Method: "GET", Path: "/customer/"},
		}

		reqPath := c.Path()
		reqMethod := c.Request().Method
		for _, condition := range skipConditions {
			if reqMethod == condition.Method && strings.HasPrefix(reqPath, condition.Path) {
				return next(c) // Skip the middleware
			}
		}

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
		}

		claims, err := authentication.VerifyJWT(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		c.Set("userID", claims.Id)
		c.Set("userEmail", claims.Email)
		c.Set("userFirstName", claims.FirstName)
		c.Set("userLastName", claims.LastName)

		return next(c)
	}
}
