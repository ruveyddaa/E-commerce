package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandlePremiumPrice(c echo.Context) error {
	return calculatePrice(c, "premium")
}

func HandleNonPremiumPrice(c echo.Context) error {
	return calculatePrice(c, "non-premium")
}

func calculatePrice(c echo.Context, role string) error {
	productID := c.QueryParam("productId")
	if productID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing productId")
	}

	basePrice := 100.0
	discount := 0.0
	if role == "premium" {
		discount = 20
	}

	finalPrice := basePrice - discount

	return c.JSON(http.StatusOK, map[string]interface{}{
		"productId":  productID,
		"role":       role,
		"basePrice":  basePrice,
		"finalPrice": finalPrice,
	})
}
