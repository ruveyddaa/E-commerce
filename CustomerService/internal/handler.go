package internal

import (
	"errors"
	"fmt"
	"tesodev-korpes/pkg/errorPackage"

	"net/http"
	"strconv"

	"tesodev-korpes/CustomerService/internal/types"

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
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceCustomerCode400101])
	}

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return errorPackage.NotFound(errorPackage.NotFoundMessages[errorPackage.ResourceCustomerCode404101])
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceCustomerCode500101])
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
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceCustomerCode400102])
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

			return pkg.ValidationFailed(details, errorPackage.ValidationErrorMessages[errorPackage.ResourceCustomerCode422101])
		}
	}

	createdID, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceCustomerCode500101])
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
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceCustomerCode400101])
	}
	var req types.UpdateCustomerRequestModel
	if err := c.Bind(&req); err != nil {
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceCustomerCode400102])
	}

	updatedCustomer, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceCustomerCode500101])
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
		return errorPackage.BadRequest(errorPackage.BadRequestMessages[errorPackage.ResourceCustomerCode400101])
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return errorPackage.NotFound(errorPackage.NotFoundMessages[errorPackage.ResourceCustomerCode404101])
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
			return errorPackage.NotFound(errorPackage.NotFoundMessages[errorPackage.ResourceCustomerCode404101])
		}

		return errorPackage.Internal(err, errorPackage.InternalServerErrorMessages[errorPackage.ResourceCustomerCode500101])
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": customers})
}
