package middleware

/*
import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
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
		}

		// Check if the current request should be skipped
		reqPath := c.Path()
		reqMethod := c.Request().Method
		for _, condition := range skipConditions {
			if reqMethod == condition.Method && strings.HasPrefix(reqPath, condition.Path) {
				return next(c) // Skip the middleware
			}
		}
		authHeader := c.Request().Header.Get("Authentication")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authentication.VerifyJWT(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		// Kullanıcı ID'sini context'e yerleştir (isteğe bağlı)
		c.Set("userID", claims.Id)

		return next(c)
	}
}
*/
