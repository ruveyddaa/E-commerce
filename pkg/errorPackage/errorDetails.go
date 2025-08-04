package errorPackage

const (
	CodeUnknown                = "UNKNOWN"
	CodeInvalidInput           = "INVALID_INPUT"
	CodeValidation             = "VALIDATION_FAILED"
	CodeNotFound               = "NOT_FOUND"
	CodeUnauthorized           = "UNAUTHORIZED"
	CodeForbidden              = "FORBIDDEN"
	CodeConflict               = "CONFLICT"
	CodeInternalError          = "INTERNAL_ERROR"
	CodeServiceDown            = "SERVICE_UNAVAILABLE"
	CodeInternalFrameworkError = "INTERNAL_FRAMEWORK_ERROR"
	CodeOrderStateConflict     = "ORDER_STATE_CONFLICT"
)

const (
	ResourceCustomer = "customer"
	ResourceOrder    = "order"
)

const (
	//-----404-----
	ResourceCustomerCode404101 = "404101"
	ResourceOrderCode404201    = "404201"
	//-----400-----
	ResourceCustomerCode400101 = "400101"
	ResourceCustomerCode400102 = "400102"

	ResourceOrderCode400201 = "400201"
	ResourceOrderCode400202 = "400202"
	//-----422-----
	ResourceCustomerCode422101 = "422101"
	//-----500-----
	ResourceCustomerCode500101  = "500101"
	ResourceOrderCode500201     = "500201"
	ResourceServiceCode500301   = "500301"
	ResourceFrameworkCode500401 = "500401"
)

var NotFoundMessages = map[string]string{
	ResourceCustomerCode404101: "The requested customer was not found.",
	ResourceOrderCode404201:    "The requested order was not found.",
}

var BadRequestMessages = map[string]string{
	ResourceCustomerCode400101: "The requested invalid customer id",
	ResourceCustomerCode400102: "The requested invalid customer body json",

	ResourceOrderCode400201: "The requested invalid order id",
	ResourceOrderCode400202: "The requested invalid order body json",
}

var ValidationErrorMessages = map[string]string{
	ResourceCustomerCode422101: "The requested customer validation error",
}

var InternalServerErrorMessages = map[string]string{
	ResourceCustomerCode500101:  "Customer service error",
	ResourceOrderCode500201:     "Order service error",
	ResourceServiceCode500301:   "Unnown service error",
	ResourceFrameworkCode500401: "A framework-related error occurred.",
}
