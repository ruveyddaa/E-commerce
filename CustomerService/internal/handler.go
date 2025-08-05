package internal

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/CustomerService/validatorCustom"
	"tesodev-korpes/pkg"

	"github.com/labstack/gommon/log"

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

	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID)
	g.GET("/email/:email", handler.GetByEmail)
	g.POST("/", handler.Create)
	g.PUT("/:id", handler.Update)
	g.DELETE("/:id", handler.Delete)
	g.GET("/list", handler.GetListCustomer)

	e.POST("/login", handler.Login)
	e.GET("/verify", handler.VerifyAuthentication)
}

func (h *Handler) Login(c echo.Context) error {
	var req types.LoginRequestModel

	if err := c.Bind(&req); err != nil {
		log.Error("Bind error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var details []pkg.ValidationErrorDetail
			for _, e := range validationErrs {
				details = append(details, pkg.ValidationErrorDetail{
					Rule:    e.Tag(),
					Message: fmt.Sprintf("The '%s' field failed on the '%s' validation", e.Field(), e.Tag()),
				})
			}
			return pkg.ValidationFailed(details, pkg.ValidationErrorMessages[pkg.ResourceCustomerCode422101])
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	email := req.Email
	log.Info("Login attempt for email: ", email)

	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email is required"})
	}

	user, err := h.service.GetByEmail(c.Request().Context(), email)
	if err != nil {
		log.Error("GetByEmail error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if user == nil {
		log.Warn("User not found for email: ", email)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	ok, err := authentication.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		log.Error("Password check failed: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Password check failed"})
	}
	if !ok {
		log.Warn("Password mismatch for email: ", email)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Generate JWT token
	/*token, err := authentication.CreateJWT(user.Id, user.FirstName, user.LastName, user.Email)
	if err != nil {
		log.Error("Token generation failed: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}
	*/
	// Convert to response model
	userResponse := ToCustomerResponse(user)

	response := types.LoginResponse{
		//Token:   token,
		User:    userResponse,
		Message: "Login successful",
	}

	log.Info("Login successful for email: ", email)
	return c.JSON(http.StatusOK, response)
}

// burdaki authorization token dogrulaması için authorizationun kendisyle alakası yok
func (h *Handler) VerifyAuthentication(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
	}

	tokenString := c.Request().Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := authentication.VerifyJWT(tokenString)
	if err != nil {
		log.Error("Token verification failed: ", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
	}

	// Verify user still exists in database
	user, err := h.service.GetByEmail(c.Request().Context(), claims.Email)
	if err != nil {
		log.Error("Error checking user existence: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	if user == nil {
		log.Error("User does not exist in database")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
	}

	response := map[string]interface{}{
		"message": "Token verified successfully",
		"user": map[string]interface{}{
			"id":         claims.Id,
			"email":      claims.Email,
			"first_name": claims.FirstName,
			"last_name":  claims.LastName,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetByEmail(c echo.Context) error {
	correlationID, _ := c.Get("CorrelationID").(string)
	email := c.Param("email")

	if !validatorCustom.IsValidEmail(email) {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400101])
	}

	customer, err := h.service.GetByEmail(c.Request().Context(), email)
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
	fmt.Println("create handler custom")

	if err := c.Bind(&req); err != nil {
		return pkg.BadRequest(pkg.BadRequestMessages[pkg.ResourceCustomerCode400102])
	}

	err := h.validate.Struct(req)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var details []pkg.ValidationErrorDetail
			fmt.Println("valide içi customer")

			for _, e := range validationErrs {
				details = append(details, pkg.ValidationErrorDetail{
					Rule:    e.Tag(),
					Message: fmt.Sprintf("The '%s' field failed on the '%s", e.Field(), e.Tag()),
				})
			}

			return pkg.ValidationFailed(details, pkg.ValidationErrorMessages[pkg.ResourceCustomerCode422101])
		}
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
