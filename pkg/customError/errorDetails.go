package customError

// ErrorDetails, uygulamanın her bir spesifik hatası için
// gerekli tüm bilgileri içeren yapıdır.
type ErrorDetails struct {
	TypeCode   string
	StatusCode int
	Message    string
	Error      error // Orijinal hatayı sarmalamak için (opsiyonel)
}

// ErrorTypeCode, uygulamadaki tüm hata kodlarını merkezi bir yapıda tutar.
// Hata yönetiminin tek ve güvenilir kaynağı (Single Source of Truth) burasıdır.
var ErrorTypeCode = map[string]ErrorDetails{
	// ----- 400 Bad Request -----
	"400101": {
		TypeCode:   "400101",
		StatusCode: 400,
		Message:    "The requested invalid customer id",
	},
	"400102": {
		TypeCode:   "400102",
		StatusCode: 400,
		Message:    "The requested invalid customer body json",
	},
	"400201": {
		TypeCode:   "400201",
		StatusCode: 400,
		Message:    "The requested invalid order id",
	},
	"400202": {
		TypeCode:   "400202",
		StatusCode: 400,
		Message:    "The requested invalid order body json",
	},

	// ----- 401 Unauthorized  -----
	"401001": {
		TypeCode:   "401001",
		StatusCode: 401,
		Message:    "Invalid email or password.",
	},
	"401002": {
		TypeCode:   "401002",
		StatusCode: 401,
		Message:    "Invalid or missing authorization token.",
	},

	// ----- 403 Forbidden  -----
	"403001": {
		TypeCode:   "403001",
		StatusCode: 403,
		Message:    "You do not have permission to access this resource.",
	},

	// ----- 404 Not Found -----
	"404101": {
		TypeCode:   "404101",
		StatusCode: 404,
		Message:    "The requested customer was not found.",
	},
	"404201": {
		TypeCode:   "404201",
		StatusCode: 404,
		Message:    "The requested order was not found.",
	},

	// ----- 409 Conflict -----
	"409201": {
		TypeCode:   "409201",
		StatusCode: 409,
		Message:    "Cannot %s order while it is in '%s' status.",
	},

	// ----- 422 Unprocessable Entity -----
	"422101": {
		TypeCode:   "422101",
		StatusCode: 422,
		Message:    "ivalid data form", //  bunlar değiştirilecek coğaltılıcak
	},

	// ----- 500 Internal Server Error -----
	"500001": { // Genel iç hata (yeni eklendi).
		TypeCode:   "500001",
		StatusCode: 500,
		Message:    "An unexpected internal error occurred.",
	},
	"500101": {
		TypeCode:   "500101",
		StatusCode: 500,
		Message:    "Customer service error",
	},
	"500201": {
		TypeCode:   "500201",
		StatusCode: 500,
		Message:    "Order service error",
	},
	"500301": {
		TypeCode:   "500301",
		StatusCode: 500,
		Message:    "Unknown service error", // 'Unnown' düzeltildi.
	},
	"500401": {
		TypeCode:   "500401",
		StatusCode: 500,
		Message:    "A framework-related error occurred.",
	},
}
