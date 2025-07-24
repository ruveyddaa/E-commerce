package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		start := time.Now()

		err := next(c)

		stop := time.Now()
		latency := stop.Sub(start)

		method := c.Request().Method
		path := c.Request().URL.Path
		status := c.Response().Status
		ip := c.RealIP()

		logrus.WithFields(logrus.Fields{
			"method":  method,
			"path":    path,
			"status":  status,
			"latency": latency,
			"ip":      ip,
		}).Info("HTTP isteği loglandı")

		return err
	}
}
