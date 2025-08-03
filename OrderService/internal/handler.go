package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg"
	"tesodev-korpes/shared/config"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title Order Service API
// @version 1.0
// @description API for managing order data
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}
	g := e.Group("/order")
	//g.POST("", handler.Create) // ← düzelt!
	g.GET("/:id", handler.GetByID)
	g.DELETE("/cancel/:id", handler.CancelOrder)
	g.PUT("/:id/ship", handler.ShipOrder)
	g.PUT("/:id/deliver", handler.DeliverOrder)
	g.GET("/list", handler.GetAllOrders)
	//g.PATCH("/cancel/:id", handler.CancelOrder)

}

func (h *Handler) GetByID(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if isValid := pkg.IsValidUUID(id); !isValid {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceOrderCode404201])
	}

	order, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return pkg.NotFound(pkg.NotFoundMessages[pkg.ResourceOrderCode404201])
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceOrderCode500201])
	}

	customer, err := fetchCustomerByID(order.CustomerId)
	if err != nil {
		pkg.LogErrorWithCorrelation(err, correlationID)
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceOrderCode500201])
	}

	response := ToOrderWithCustomerResponse(order, customer)

	pkg.LogInfoWithCorrelation("Order with customer fetched", correlationID)
	return c.JSON(http.StatusOK, response)
}

// buranın errorlarını düzeltelim
func fetchCustomerByID(customerID string) (*types.CustomerResponseModel, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customerID is empty")
	}

	baseURL := config.GetServiceURLs().CustomerServiceURL

	url := fmt.Sprintf("%s/customer/%s", baseURL, customerID)

	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("customer not found, status: %d", resp.StatusCode)
	}

	var customer types.CustomerResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (h *Handler) CancelOrder(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if isValid := pkg.IsValidUUID(id); !isValid {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceOrderCode404201])
	}

	err := h.service.CancelOrder(c.Request().Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("order not found for ID: %s", id) {
			return pkg.NotFound(pkg.NotFoundMessages[pkg.ResourceOrderCode404201])
		}

		if errResp, ok := err.(*pkg.AppError); ok && errResp.Code == pkg.CodeOrderStateConflict {
			return c.JSON(http.StatusConflict, echo.Map{"message": errResp.Message})
		} // buranın errorlarınını düzeltelim

		pkg.LogErrorWithCorrelation(err, correlationID)
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceOrderCode500201])
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order cancelled successfully. The order is now inactive."})
}

func (h *Handler) ShipOrder(c echo.Context) error {
	id := c.Param("id")

	err := h.service.ShipOrder(c.Request().Context(), id)
	if err != nil {
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceOrderCode500201])
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order shipped successfully"})
}

// DeliverOrder godoc
// @Summary Deliver an order
// @Description Changes the order status to DELIVERED if current status is SHIPPED
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order delivered successfully"
// @Failure 400 {object} pkg.AppError "Invalid ID format"
// @Failure 404 {object} pkg.AppError "Order not found"
// @Failure 409 {object} pkg.AppError "Invalid order state for delivery"
// @Failure 500 {object} pkg.AppError "Internal server error"
// @Security ApiKeyAuth
// @Router /order/{id}/deliver [put]
func (h *Handler) DeliverOrder(c echo.Context) error {
	id := c.Param("id")
	if !pkg.IsValidUUID(id) {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceOrderCode404201])
	}

	err := h.service.DeliverOrder(c.Request().Context(), id)
	if err != nil {
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceOrderCode500201])
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order delivered successfully"})
}

func (h *Handler) GetAllOrders(c echo.Context) error {
	params := types.Pagination{
		Page:  1,
		Limit: 10,
	}

	if p := c.QueryParam("page"); p != "" {
		if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
			params.Page = pageInt
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if limitInt, err := strconv.Atoi(l); err == nil && limitInt > 0 {
			params.Limit = limitInt
		}
	}

	orders, err := h.service.GetAllOrders(c.Request().Context(), params)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return pkg.NotFound(pkg.NotFoundMessages[pkg.ResourceOrderCode404201])
		}

		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceOrderCode500201])
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": orders})

}
