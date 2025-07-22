package internal

import (
	"net/http"
	"strconv"
	"tesodev-korpes/CustomerService/internal/types"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}

	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID)
	g.POST("/", handler.Create)
	g.PUT("/:id", handler.Update)
	g.DELETE("/:id", handler.Delete)
	g.GET("/list", handler.GetListCustomer)
}

func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id") // URL’den gelen id’yi alıyoruz

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Customer not found"}) // TODO: rüveyda errorları düzeltecek
	}

	return c.JSON(http.StatusOK, customer) // frontend'e JSON olarak döner
}

func (h *Handler) Create(c echo.Context) error {
	var req types.CreateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	createdID, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message":   "Succeeded!",
		"createdId": createdID,
	})
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")

	var req types.UpdateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return Respond(c, NewBadRequest("Invalid request body"), "Failed to bind request")
	}

	// TODO:  tüm erroları düzeltip ortak işleyiş belirlenicek

	updatedCustomer, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return Respond(c, err, "Failed to update customer")
	}

	response := ToCustomerResponse(updatedCustomer)
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetListCustomer(c echo.Context) error {
	params := types.Pagination{
		Limit: 10,
		Page:  1,
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

	customers, err := h.service.Get(c.Request().Context(), params)
	if err != nil {
		return NewInternal(err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": customers})
}
