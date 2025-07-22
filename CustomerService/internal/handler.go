package internal

import (
	"net/http"
	"strconv"
	"strings"
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

// GetByID godoc
// @Summary Get customer by ID
// @Description Retrieve a customer by its ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 404 {object} map[string]string
// @Router /customer/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {

		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, customer)
}

// Create godoc
// @Summary Create a new customer
// @Description Create customer from payload
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body types.CreateCustomerRequestModel true "Customer payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /customer/ [post]
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

// Update godoc
// @Summary Update customer
// @Description Update customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body types.UpdateCustomerRequestModel true "Customer payload"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 400,500 {object} string
// @Router /customer/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")

	var req types.UpdateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	updatedCustomer, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := ToCustomerResponse(updatedCustomer)
	return c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Delete customer
// @Description Delete a customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 204 "No Content"
// @Failure 400,404 {object} map[string]string
// @Router /customer/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return Respond(c, err, "Failed to delete customer")
	}
	return c.NoContent(http.StatusNoContent)
}

// GetListCustomer godoc
// @Summary List customers with pagination
// @Description Retrieve a paginated list of customers
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /customer/list [get]
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
