package customError

import "net/http"

type ErrorDetails struct {
	TypeCode   int
	StatusCode int
	Message    string
	Error      error
}

type ErrorKey string

var ErrorDefinitions = map[ErrorKey]ErrorDetails{
	// ----- 400 Bad Request -----
	InvalidCustomerID: {
		TypeCode:   400101,
		StatusCode: http.StatusBadRequest,
		Message:    "The requested invalid customer id.",
	},
	InvalidCustomerBody: {
		TypeCode:   400102,
		StatusCode: http.StatusBadRequest,
		Message:    "The requested invalid customer body json.",
	},
	InvalidOrderID: {
		TypeCode:   400201,
		StatusCode: http.StatusBadRequest,
		Message:    "The requested invalid order id.",
	},
	InvalidOrderBody: {
		TypeCode:   400202,
		StatusCode: http.StatusBadRequest,
		Message:    "The requested invalid order body json.",
	},

	// ----- 401 Unauthorized -----
	InvalidCredentials: {
		TypeCode:   401001,
		StatusCode: http.StatusUnauthorized,
		Message:    "Invalid email or password.",
	},
	MissingAuthToken: {
		TypeCode:   401002,
		StatusCode: http.StatusUnauthorized,
		Message:    "Invalid or missing authorization token.",
	},

	// ----- 403 Forbidden -----
	ForbiddenAccess: {
		TypeCode:   403001,
		StatusCode: http.StatusForbidden,
		Message:    "You do not have permission to access this resource.",
	},

	// ----- 404 Not Found -----
	CustomerNotFound: {
		TypeCode:   404101,
		StatusCode: http.StatusNotFound,
		Message:    "The requested customer was not found.",
	},
	OrderNotFound: {
		TypeCode:   404201,
		StatusCode: http.StatusNotFound,
		Message:    "The requested order was not found.",
	},

	// ----- 409 Conflict -----
	OrderStatusConflict: {
		TypeCode:   409201,
		StatusCode: http.StatusConflict,
		Message:    "Cannot %s order while it is in '%s' status.",
	},

	// ----- 422 Unprocessable Entity -----
	InvalidDataFormat: {
		TypeCode:   422101,
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "Invalid data format provided.",
	},

	// ----- 500 Internal Server Error -----
	InternalServerError: {
		TypeCode:   500001,
		StatusCode: http.StatusInternalServerError,
		Message:    "An unexpected internal error occurred.",
	},
	CustomerServiceError: {
		TypeCode:   500101,
		StatusCode: http.StatusInternalServerError,
		Message:    "An error occurred in the customer service.",
	},
	OrderServiceError: {
		TypeCode:   500201,
		StatusCode: http.StatusInternalServerError,
		Message:    "An error occurred in the order service.",
	},
	UnknownServiceError: {
		TypeCode:   500301,
		StatusCode: http.StatusInternalServerError,
		Message:    "An unknown service error occurred.",
	},
	FrameworkError: {
		TypeCode:   500401,
		StatusCode: http.StatusInternalServerError,
		Message:    "A framework-related error occurred.",
	},
}

// Hata anahtarları (constants)
const (
	// Bad Request
	InvalidCustomerID   ErrorKey = "InvalidCustomerID"
	InvalidCustomerBody ErrorKey = "InvalidCustomerBody"
	InvalidOrderID      ErrorKey = "InvalidOrderID"
	InvalidOrderBody    ErrorKey = "InvalidOrderBody"

	// Unauthorized
	InvalidCredentials ErrorKey = "InvalidCredentials"
	MissingAuthToken   ErrorKey = "MissingAuthToken"

	// Forbidden
	ForbiddenAccess ErrorKey = "ForbiddenAccess"

	// Not Found
	CustomerNotFound ErrorKey = "CustomerNotFound"
	OrderNotFound    ErrorKey = "OrderNotFound"

	// Conflict
	OrderStatusConflict ErrorKey = "OrderStatusConflict"

	// Unprocessable Entity
	InvalidDataFormat ErrorKey = "InvalidDataFormat"

	// Internal Server Error
	InternalServerError  ErrorKey = "InternalServerError"
	CustomerServiceError ErrorKey = "CustomerServiceError"
	OrderServiceError    ErrorKey = "OrderServiceError"
	UnknownServiceError  ErrorKey = "UnknownServiceError"
	FrameworkError       ErrorKey = "FrameworkError"
)
