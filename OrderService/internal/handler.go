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
	correlationID, _ := c.Get("CorrelationID").(string)
	userID, ok := c.Get("userID").(string)
	if !ok {
		return errorPackage.Unauthorized()
	}
	userEmail, ok := c.Get("userEmail").(string)
	if !ok {
		return errorPackage.Unauthorized()
	}
	var req types.CreateOrderRequestModel
	if err := c.Bind(&req); err != nil {
		return errorPackage.BadRequest("Invalid request data: " + err.Error())
	}
	req.CustomerId = userID

	if err := h.validate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var details []pkg.ValidationErrorDetail
			for _, e := range validationErrs {
				details = append(details, pkg.ValidationErrorDetail{
					Rule:    e.Tag(),
					Message: fmt.Sprintf("The '%s' field failed on the '%s' rule", e.Field(), e.Tag()),
				})
			}
			return pkg.ValidationFailed(details, errorPackage.ValidationErrorMessages[errorPackage.ResourceCustomerCode422101])
		}
		return errorPackage.BadRequest("Validation error")
	}
	customer, err := h.service.fetchCustomerByID(userID)
	if err != nil {
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.Internal(err, "Customer service connection failed")
	}
	if customer == nil {
		return errorPackage.NotFound(errorPackage.NotFoundMessages[errorPackage.ResourceCustomerCode404101])
	}
	order := FromCreateOrderRequest(&req)
	createdID, err := h.service.Create(c.Request().Context(), order)
	if err != nil {
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceOrderCode500201])
	}
	createdOrder, err := h.service.GetByID(c.Request().Context(), createdID)
	if err != nil {
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.Internal(err, "Failed to retrieve the created order")
	}
	pkg.LogInfoWithCorrelation("Order created successfully", correlationID)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"order": createdOrder,
		"user": map[string]interface{}{
			"id":    userID,
			"email": userEmail,
		},
		"message": "Order created successfully",
	})
}

func (h *Handler) GetByID(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")

	if !pkg.IsValidUUID(id) {
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceOrderCode404201])
	}

	orderWithCustomer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errorPackage.NotFound(errorPackage.NotFoundMessages[errorPackage.ResourceOrderCode404201])
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceOrderCode500201])
	}

	pkg.LogInfoWithCorrelation("Order with customer fetched", correlationID)
	return c.JSON(http.StatusOK, orderWithCustomer)
}

func (h *Handler) ShipOrder(c echo.Context) error {
	id := c.Param("id")

	err := h.service.ShipOrder(c.Request().Context(), id)
	if err != nil {
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceOrderCode500201])
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
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceOrderCode404201])
	}

	err := h.service.DeliverOrder(c.Request().Context(), id)
	if err != nil {
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceOrderCode500201])
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order delivered successfully"})
}

func (h *Handler) CancelOrder(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if isValid := pkg.IsValidUUID(id); !isValid {
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceOrderCode404201])
	}

	err := h.service.CancelOrder(c.Request().Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("order not found for ID: %s", id) {
			return errorPackage.NotFound(errorPackage.NotFoundMessages[errorPackage.ResourceOrderCode404201])
		}

		if errResp, ok := err.(*errorPackage.AppError); ok && errResp.Code == errorPackage.CodeOrderStateConflict {
			return c.JSON(http.StatusConflict, echo.Map{"message": errResp.Message})
		}

		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceOrderCode500201])
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
		if errResp, ok := err.(*pkg.AppError); ok && errResp.Code == pkg.CodeOrderStateConflict {
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
			return errorPackage.NotFound(errorPackage.NotFoundMessages[errorPackage.ResourceOrderCode404201])
		}

		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceOrderCode500201])
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": orders})

}
