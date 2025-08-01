package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		start := time.Now()
		method := c.Request().Method
		path := c.Request().URL.Path
		err := next(c)

		stop := time.Now()

		correlationID, _ := c.Get("CorrelationID").(string)

		latency := stop.Sub(start)

		status := c.Response().Status
		ip := c.RealIP()

		logrus.WithFields(logrus.Fields{
			"method":        method,
			"path":          path,
			"status":        status,
			"latency":       latency,
			"ip":            ip,
			"correlationID": correlationID,
		}).Info("HTTP isteği loglandı")

		return err
	}
}
