package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Skipper(c echo.Context) bool {
	path := c.Path()
	method := c.Request().Method

	if path == "/customer/login" && method == http.MethodPost {
		return true
	}
	if path == "/customer/create" && method == http.MethodPost {
		return true
	}
	return false
}
