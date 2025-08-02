package internal

import (
	"errors"

	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/CustomerService/validator"
	"tesodev-korpes/pkg"
	_ "tesodev-korpes/pkg/middleware"

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

	e.POST("/login", handler.Login)
	//e.GET("/verify", handler.Verify)
}
func (h *Handler) Login(c echo.Context) error {
	var req types.LoginRequestModel

	if err := c.Bind(&req); err != nil {
		log.Error("Bind error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	email := req.Email
	log.Info("Extracted email: ", email)

	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email is required"})
	}

	user, err := h.service.GetByEmail(c.Request().Context(), email)
	if err != nil {
		log.Error("GetByEmail error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if user == nil {
		log.Warn("User not found")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	ok, err := authentication.CheckPasswordHash(req.Password, user.Password)

	if err != nil {
		log.Error("Password check failed: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Password check failed"})
	}
	if !ok {
		log.Warn("Password mismatch")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	log.Info("Login successful")
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetByEmail(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("email")

	if !validator.IsValidEmail(id) { // Bu fonksiyonu senin yazman gerek
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400101])
	}

	customer, err := h.service.GetByEmail(c.Request().Context(), id)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return pkg.NotFound(pkg.NotFoundMessages[pkg.ResourceCustomerCode404101])
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceCustomerCode500101])
	}
	pkg.LogInfoWithCorrelation("Customer found", correlationID)
	return c.JSON(http.StatusOK, customer)
}

// GetByID godoc
// @Summary Get customer by ID
// @Description Get a customer by its unique ID
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 400 {object} pkg.AppError "Invalid ID format"
// @Failure 404 {object} pkg.AppError "Customer not found"
// @Failure 500 {object} pkg.AppError "Internal server error"
// @Router /customer/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if isValidID := pkg.IsValidUUID(id); !isValidID {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400101])
	}

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return pkg.NotFound(pkg.NotFoundMessages[pkg.ResourceCustomerCode404101])
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceCustomerCode500101])
	}
	pkg.LogInfoWithCorrelation("Customer found", correlationID)
	return c.JSON(http.StatusOK, customer)
}

// Create godoc
// @Summary Create a new customer
// @Description Create a new customer with the given data
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param customer body types.CreateCustomerRequestModel true "Customer to create"
// @Success 201 {object} map[string]interface{} "Returns created customer ID"
// @Failure 400 {object} pkg.AppError "Invalid request body"
// @Failure 500 {object} pkg.AppError "Internal server error"
// @Router /customer/ [post]
func (h *Handler) Create(c echo.Context) error {
	var req types.CreateCustomerRequestModel

	if err := c.Bind(&req); err != nil {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400102])
	}

	createdID, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceCustomerCode500101])
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
// @Security ApiKeyAuth
// @Param id path string true "Customer ID"
// @Param customer body types.UpdateCustomerRequestModel true "Customer data to update"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 400 {object} pkg.AppError "Invalid ID format or request body"
// @Failure 404 {object} pkg.AppError "Customer not found"
// @Failure 500 {object} pkg.AppError "Internal server error"
// @Router /customer/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	if isValidID := pkg.IsValidUUID(id); !isValidID {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400101])
	}
	var req types.UpdateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400102])
	}

	updatedCustomer, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceCustomerCode500101])
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
// @Security ApiKeyAuth
// @Param id path string true "Customer ID"
// @Success 204 "No Content"
// @Failure 400 {object} pkg.AppError "Invalid ID format"
// @Failure 404 {object} pkg.AppError "Customer not found"
// @Failure 500 {object} pkg.AppError "Internal server error"
// @Router /customer/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if isValidID := pkg.IsValidUUID(id); !isValidID {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400101])
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return pkg.NotFound(pkg.NotFoundMessages[pkg.ResourceCustomerCode404101])
	}
	return c.NoContent(http.StatusNoContent)
}

// GetListCustomer godoc
// @Summary List customers with pagination
// @Description Retrieve a paginated list of customers
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page"
// @Success 200 {object} map[string]interface{} "Returns list of customers"
// @Failure 500 {object} pkg.AppError "Internal server error"
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
			return pkg.NotFound(pkg.NotFoundMessages[pkg.ResourceCustomerCode404101])
		}

		return pkg.Internal(err, pkg.InternalServerErrorMessages[pkg.ResourceCustomerCode500101])
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": customers})
}
