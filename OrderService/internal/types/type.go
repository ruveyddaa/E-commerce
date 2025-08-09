package types

import "time"

type CreateOrderRequestModel struct {
	CustomerId      string      `json:"customer_id,omitempty" validate:"required,dive,required"`
	Items           []OrderItem `json:"items" validate:"required,dive,required"`
	ShippingAddress Address     `json:"shipping_address" validate:"required"`
	BillingAddress  Address     `json:"billing_address" validate:"required"`
}

type OrderResponseModel struct {
	Id              string      `json:"id"`
	CustomerId      string      `json:"customer_id"`
	Items           []OrderItem `json:"items"`
	ShippingAddress Address     `json:"shipping_address"`
	BillingAddress  Address     `json:"billing_address"`
	TotalPrice      float64     `json:"total_price"`
	Status          string      `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	IsDelete        bool        `json:"is_delete"`
}
type CustomerResponseModel struct {
	Id        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     []Phone   `json:"phone"`
	Address   []Address `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Pagination struct {
	Page  int
	Limit int
}

type OrderWithCustomerResponse struct {
	OrderResponseModel
	Customer CustomerResponseModel `json:"customer,omitempty"`
}
