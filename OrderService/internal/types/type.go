package types

import "time"

type CreateOrderRequestModel struct {
	CustomerId      string      `json:"customer_id" binding:"required,uuid4"`
	Items           []OrderItem `json:"items" binding:"required,dive,required"`
	ShippingAddress Address     `json:"shipping_address" binding:"required"`
	BillingAddress  Address     `json:"billing_address" binding:"required"`
}

type OrderResponseModel struct {
	Id              string      `json:"id"`
	CustomerId      string      `json:"customer_id"`
	Items           []OrderItem `json:"items"`
	ShippingAddress Address     `json:"shipping_address"`
	BillingAddress  Address     `json:"billing_address"`
	TotalPrice      float64     `json:"total_price"`
	Status          OrderStatus `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	IsActive        bool        `json:"is_active"`
}

type Pagination struct {
	Page  int
	Limit int
}

type OrderWithCustomerResponse struct {
	OrderResponseModel
	Customer interface{} `json:"customer,omitempty"`
}
