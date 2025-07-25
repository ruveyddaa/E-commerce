package internal

import (
	"errors"
	"net/http"
	"strconv"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/pkg"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title Customer Service API
// @version 1.0
// @description API for managing customer data
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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
// @Description Get a customer by its unique ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /customer/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return pkg.NotFound()
		}

		if err.Error() == "the provided hex string is not a valid ObjectID" {
			return pkg.BadRequest(err.Error())
		}

		return pkg.Internal(err)
	}

	return c.JSON(http.StatusOK, customer)
}

// Create godoc
// @Summary Create a new customer
// @Description Create a new customer with the given data
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body types.CreateCustomerRequestModel true "Customer to create"
// @Success 201 {object} map[string]interface{} "Returns created customer ID"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /customer/ [post]
func (h *Handler) Create(c echo.Context) error {
	var req types.CreateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return pkg.BadRequest(err.Error())
	}

	createdID, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return pkg.Internal(err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message":   "Succeeded!",
		"createdId": createdID,
	})
}

// Update godoc
// @Summary Update an existing customer
// @Description Update a customer with the given ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body types.UpdateCustomerRequestModel true "Customer data to update"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 400 {object} map[string]string "Invalid ID format or request body"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /customer/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")

	var req types.UpdateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return pkg.BadRequest(err.Error())
	}

	updatedCustomer, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return pkg.Internal(err)
	}


	response := ToCustomerResponse(updatedCustomer)
	return c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Delete a customer by ID
// @Description Delete a customer from the system
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /customer/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return pkg.NotFound()
	}
	return c.NoContent(http.StatusNoContent)
}

// GetListCustomer godoc
// @Summary List customers with pagination
// @Description Retrieve a paginated list of customers
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page"
// @Success 200 {object} map[string]interface{} "Returns list of customers"
// @Failure 500 {object} map[string]string "Internal server error"
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return pkg.NotFound()
		}

		return pkg.Internal(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": customers})
}
