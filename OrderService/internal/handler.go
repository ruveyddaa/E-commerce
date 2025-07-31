package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}
	g := e.Group("/order")
	g.POST("", handler.Create) // ← düzelt!
	g.GET("/:id", handler.GetByID)
	g.DELETE("/:id", handler.Delete)
}

// Order + Customer JSON response modeli
type OrderWithCustomerResponse struct {
	types.OrderResponseModel
	Customer interface{} `json:"customer,omitempty"`
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

	response := OrderWithCustomerResponse{
		OrderResponseModel: *order,
		Customer:           customer,
	}

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

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := h.service.DeleteOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Order not found",
		})
	}

	return c.NoContent(http.StatusNoContent)
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
