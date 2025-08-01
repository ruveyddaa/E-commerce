package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg"
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
	g.POST("", handler.Create) // ← düzelt!
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

	order, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return pkg.NotFound()
		}
		if err.Error() == "the provided hex string is not a valid ObjectID" {
			return pkg.BadRequest(err.Error())
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return pkg.Internal(err)
	}

	customer, err := fetchCustomerByID(order.CustomerId)
	if err != nil {
		pkg.LogErrorWithCorrelation(err, correlationID)
		return pkg.Internal(errors.New("failed to fetch customer info"))
	}

	response := ToOrderWithCustomerResponse(order, customer)

	pkg.LogInfoWithCorrelation("Order with customer fetched", correlationID)
	return c.JSON(http.StatusOK, response)
}

// Customer API'ye HTTP GET atan yardımcı fonksiyon
func fetchCustomerByID(customerID string) (interface{}, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customerID boş")
	}

	url := fmt.Sprintf("http://localhost:8001/customer/%s", customerID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("customer not found, status: %d", resp.StatusCode)
	}

	var customer interface{}
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (h *Handler) CancelOrder(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")

	err := h.service.CancelOrder(c.Request().Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("order not found for ID: %s", id) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Order not found"})
		}

		if errResp, ok := err.(*pkg.AppError); ok && errResp.Code == pkg.CodeOrderStateConflict {
			return c.JSON(http.StatusConflict, echo.Map{"message": errResp.Message})
		}

		pkg.LogErrorWithCorrelation(err, correlationID)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order cancelled successfully. The order is now inactive."})
}

func (h *Handler) Create(c echo.Context) error {
	var req types.Order

	if err := c.Bind(&req); err != nil {
		return pkg.BadRequest("Geçersiz istek verisi: " + err.Error())
	}

	createdID, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return pkg.Internal(err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message":   "Order başarıyla oluşturuldu",
		"createdId": createdID,
	})
}

func (h *Handler) ShipOrder(c echo.Context) error {
	id := c.Param("id")

	err := h.service.ShipOrder(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
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

	err := h.service.DeliverOrder(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
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
			return pkg.NotFound()
		}

		return pkg.Internal(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": orders})

}
