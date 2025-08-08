package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tesodev-korpes/pkg/errorPackage"

	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg"

	"github.com/go-playground/validator/v10"
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
	service  *Service
	validate *validator.Validate
}

func NewHandler(e *echo.Echo, service *Service) {
	validate := validator.New()

	handler := &Handler{
		service:  service,
		validate: validate,
	}

	g := e.Group("/order")
	g.POST("", handler.Create)
	g.GET("/:id", handler.GetByID)
	g.PATCH("/:id/ship", handler.ShipOrder)
	g.PATCH("/:id/deliver", handler.DeliverOrder)
	g.PATCH("/cancel/:id", handler.CancelOrder)
	g.PATCH("/delete/:id", handler.DeleteOrder)
	g.GET("/list", handler.GetAllOrders)
}

func (h *Handler) Create(c echo.Context) error {
	var req types.CreateOrderRequestModel

	if err := c.Bind(&req); err != nil {
		return errorPackage.NewBadRequest("400102")
	}

	order := FromCreateOrderRequest(&req)

	createdID, err := h.service.Create(c.Request().Context(), order)
	if err != nil {
		return errorPackage.NewInternal("500201", err)
	}

	createdOrder, err := h.service.GetByID(c.Request().Context(), createdID)
	if err != nil {
		return errorPackage.NewInternal("500201", err)
	}

	return c.JSON(http.StatusCreated, createdOrder)
}

func (h *Handler) GetByID(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")

	if !pkg.IsValidUUID(id) {
		return errorPackage.NewBadRequest("400201")
	}

	orderWithCustomer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errorPackage.NewNotFound("404201")
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.NewInternal("500201", err)
	}

	pkg.LogInfoWithCorrelation("Order with customer fetched", correlationID)
	return c.JSON(http.StatusOK, orderWithCustomer)
}

func (h *Handler) ShipOrder(c echo.Context) error {
	id := c.Param("id")

	err := h.service.ShipOrder(c.Request().Context(), id)
	if err != nil {
		return errorPackage.NewInternal("500201", err)
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
		return errorPackage.NewBadRequest("404201")
	}

	err := h.service.DeliverOrder(c.Request().Context(), id)
	if err != nil {
		return errorPackage.NewInternal("500201", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order delivered successfully"})
}

func (h *Handler) CancelOrder(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if isValid := pkg.IsValidUUID(id); !isValid {
		return errorPackage.NewBadRequest("400201")
	}

	err := h.service.CancelOrder(c.Request().Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("order not found for ID: %s", id) {
			return errorPackage.NewNotFound("404201")
		}

		if errResp, ok := err.(*errorPackage.AppError); ok {
			return c.JSON(http.StatusConflict, echo.Map{"message": errResp.Message})
		}

		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.NewInternal("500201", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order cancelled successfully. "})
}

func (h *Handler) DeleteOrder(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")

	err := h.service.DeleteOrder(c.Request().Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("order not found for ID: %s", id) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Order not found"})
		}
		if errResp, ok := err.(*errorPackage.AppError); ok {
			return c.JSON(http.StatusConflict, echo.Map{"message": errResp.Message})
		}

		pkg.LogErrorWithCorrelation(err, correlationID)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order deleted (soft delete) successfully."})
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
			return errorPackage.NewNotFound("404201")
		}
		return errorPackage.NewInternal("500201", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": orders})

}
