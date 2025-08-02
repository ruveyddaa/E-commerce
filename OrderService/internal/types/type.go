package types

import "time"

type CreateOrderRequestModel struct {
	CustomerId      string      `bson:"customer_id" json:"customer_id" validate:"required,uuid4"`
	Items           []OrderItem `bson:"items" json:"items" validate:"required,dive,required"`
	ShippingAddress Address     `bson:"shipping_address" json:"shipping_address" validate:"required"`
	BillingAddress  Address     `bson:"billing_address" json:"billing_address" validate:"required"`
}

type OrderResponseModel struct {
	Id              string      `bson:"_id" json:"id"`
	CustomerId      string      `bson:"customer_id" json:"customer_id"`
	Items           []OrderItem `bson:"items" json:"items"`
	ShippingAddress Address     `bson:"shipping_address" json:"shipping_address"`
	BillingAddress  Address     `bson:"billing_address" json:"billing_address"`
	TotalPrice      float64     `bson:"total_price" json:"total_price"`
	Status          OrderStatus `bson:"status" json:"status"`
	IsActive        bool        `bson:"is_active" json:"is_active"`
	CreatedAt       time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `bson:"updated_at" json:"updated_at"`
}

type CustomerResponseModel struct {
	Id        string            `json:"id" bson:"_id,omitempty"`
	FirstName string            `json:"first_name" bson:"first_name"`
	LastName  string            `json:"last_name" bson:"last_name"`
	Email     map[string]string `json:"email" bson:"email"`
	Phone     []Phone           `json:"phone" bson:"phone"`
	Address   []Address         `json:"address" bson:"address"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
	IsActive  bool              `json:"is_active" bson:"is_active"`
}

type Pagination struct {
	Page  int
	Limit int
}

type OrderWithCustomerResponse struct {
	OrderResponseModel
	Customer interface{} `json:"customer,omitempty"`
}
