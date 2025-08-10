package internal

import (
	"errors"
	"fmt"
	"tesodev-korpes/pkg/auth"

	"tesodev-korpes/pkg/errorPackage"
	"tesodev-korpes/pkg/middleware"

	"net/http"
	"strconv"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/CustomerService/validatorCustom"
	"tesodev-korpes/pkg"

	"github.com/go-playground/validator/v10"
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
	service  *Service
	validate *validator.Validate
}

func NewHandler(e *echo.Echo, service *Service) {
	validate := validator.New()
	handler := &Handler{
		service:  service,
		validate: validate,
	}

	allowedRole_premium := []string{"premium", "non-premium"}

	e.Use(middleware.Authentication(auth.Skipper))

	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID)
	g.GET("/email/:email", handler.GetByEmail, middleware.AuthorizationMiddleware(allowedRole_premium))
	g.PUT("/:id", handler.Update)
	g.DELETE("/:id", handler.Delete)
	g.GET("/list", handler.GetListCustomer)
	g.GET("/verify", handler.VerifyAuthentication)

	g.POST("/create", handler.Create)
	g.POST("/login", handler.Login)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, return access token and user info.
// @Tags authentication
// @Accept json
// @Produce json
// @Param loginRequest body types.LoginRequestModel true "Login credentials"
// @Success 200 {object} types.LoginResponse "Returns access token and customer info"
// @Failure 400 {object} errorPackage.AppError "Invalid request payload"
// @Failure 401 {object} errorPackage.AppError "Unauthorized, invalid credentials"
// @Failure 422 {object} errorPackage.AppError "Validation error on input data"
// @Failure 500 {object} errorPackage.AppError "Internal server error"
// @Router /customer/login [post]
func (h *Handler) Login(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)

	var req types.LoginRequestModel
	if err := c.Bind(&req); err != nil {
		return errorPackage.NewBadRequest("400102")
	}

	if err := h.validate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var details []pkg.ValidationErrorDetail
			for _, e := range validationErrs {
				details = append(details, pkg.ValidationErrorDetail{
					Rule:    e.Tag(),
					Message: fmt.Sprintf("The '%s' field failed on the '%s'", e.Field(), e.Tag()),
				})
			}
			return errorPackage.NewValidation("422101")
		}
	}

	token, customer, err := h.service.Login(c.Request().Context(), req.Email, req.Password, correlationID)
	if err != nil {
		return errorPackage.NewUnauthorized("401001")
	}

	response := ToLoginResponse(token, customer)
	return c.JSON(http.StatusOK, response)
}

// VerifyAuthentication godoc
// @Summary Verify user authentication
// @Description Verify the authentication token and return user details if valid.
// @Tags authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} types.VerifyTokenResponse "Returns user info on successful authentication"
// @Failure 401 {object} errorPackage.AppError "Unauthorized or invalid token"
// @Failure 500 {object} errorPackage.AppError "Internal server error"
// @Router /customer/verify [get]
func (h *Handler) VerifyAuthentication(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)

	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return errorPackage.NewUnauthorized("401002")
	}

	user, err := h.service.GetByID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errorPackage.NewUnauthorized("401002")
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.NewInternal("500101", err)
	}

	response := ToVerifyTokenResponse(user)
	return c.JSON(http.StatusOK, response)
}

// GetByEmail godoc
// @Summary Get customer by email
// @Description Retrieve a customer by their email address. Used for login verification or profile retrieval.
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email path string true "Customer Email"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 400 {object} errorPackage.AppError "Invalid email format"
// @Failure 404 {object} errorPackage.AppError "Customer not found"
// @Failure 500 {object} errorPackage.AppError "Internal server error"
// @Router /customer/email/{email} [get]
func (h *Handler) GetByEmail(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	email := c.Param("email")

	if !validatorCustom.IsValidEmail(email) {
		return errorPackage.NewBadRequest("400102")
	}

	customer, err := h.service.GetByEmail(c.Request().Context(), email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errorPackage.NewNotFound("404101")
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.NewInternal("500101", err)
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
// @Failure 400 {object} errorPackage.AppError "Invalid ID format"
// @Failure 404 {object} errorPackage.AppError "Customer not found"
// @Failure 500 {object} errorPackage.AppError "Internal server error"
// @Router /customer/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	id := c.Param("id")
	if isValidID := pkg.IsValidUUID(id); !isValidID {
		return errorPackage.NewBadRequest("400101")
	}

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return errorPackage.NewNotFound("404101")
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.NewInternal("500101", err)
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
// @Failure 400 {object} errorPackage.AppError "Invalid request body"
// @Failure 500 {object} errorPackage.AppError "Internal server error"
// @Router /customer/ [post]
func (h *Handler) Create(c echo.Context) error {
	var req types.CreateCustomerRequestModel
	fmt.Println("create handler custom")

	if err := c.Bind(&req); err != nil {
		return errorPackage.NewBadRequest("400102")
	}

	err := h.validate.Struct(req)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var details []pkg.ValidationErrorDetail
			fmt.Println("Customer validation in progress")

			for _, e := range validationErrs {
				details = append(details, pkg.ValidationErrorDetail{
					Rule:    e.Tag(),
					Message: fmt.Sprintf("The '%s' field failed on the '%s", e.Field(), e.Tag()),
				})
			}

			return errorPackage.NewValidation("422101")
		}
	}

	createdID, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return errorPackage.NewInternal("500101", err)
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
// @Failure 400 {object} errorPackage.AppError "Invalid ID format or request body"
// @Failure 404 {object} errorPackage.AppError "Customer not found"
// @Failure 500 {object} errorPackage.AppError "Internal server error"
// @Router /customer/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	if !pkg.IsValidUUID(id) {
		return errorPackage.NewBadRequest("400101")
	}

	var req types.UpdateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return errorPackage.NewBadRequest("400102")
	}

	existing, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return errorPackage.NewInternal("500102", err)
	}

	updatedCustomer := FromUpdateCustomerRequest(existing, &req)

	if err := h.service.Update(c.Request().Context(), id, updatedCustomer); err != nil {
		return errorPackage.NewInternal("500101", err)
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
// @Failure 400 {object} errorPackage.AppError "Invalid ID format"
// @Failure 404 {object} errorPackage.AppError "Customer not found"
// @Failure 500 {object} errorPackage.AppError "Internal server error"
// @Router /customer/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if isValidID := pkg.IsValidUUID(id); !isValidID {
		return errorPackage.NewBadRequest("400101")
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return errorPackage.NewNotFound("404101")
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
// @Failure 500 {object} errorPackage.AppError "Internal server error"
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
			return errorPackage.NewNotFound("404101")
		}

		return errorPackage.NewInternal("500101", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": customers})
}
