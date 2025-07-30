package types

import "time"

type CreateOrderRequestModel struct {
	CustomerId      string      `json:"customer_id" binding:"required"`
	Items           []OrderItem `json:"items" binding:"required,dive,required"`
	ShippingAddress Address     `json:"shipping_address" binding:"required"`
	BillingAddress  Address     `json:"billing_address" binding:"required"`
}

type UpdateOrderRequestModel struct {
	Items           []OrderItem    `json:"items,omitempty"`
	ShippingAddress *Address       `json:"shipping_address,omitempty"`
	BillingAddress  *Address       `json:"billing_address,omitempty"`
	Status          *OrderStatus   `json:"status,omitempty"`
	PaymentStatus   *PaymentStatus `json:"payment_status,omitempty"`
	IsActive        *bool          `json:"is_active,omitempty"`
}

type OrderResponseModel struct {
	Id              string        `json:"id"`
	CustomerId      string        `json:"customer_id"`
	Items           []OrderItem   `json:"items"`
	ShippingAddress Address       `json:"shipping_address"`
	BillingAddress  Address       `json:"billing_address"`
	TotalPrice      float64       `json:"total_price"`
	Status          OrderStatus   `json:"status"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	IsActive        bool          `json:"is_active"`
}

type Pagination struct {
	Page  int
	Limit int
}
