package types

import "time"

type Order struct {
	Id              string      `bson:"_id,omitempty" json:"id"`
	CustomerId      string      `bson:"customer_id" json:"customer_id"`
	Items           []OrderItem `bson:"items" json:"items"`
	ShippingAddress Address     `bson:"shipping_address" json:"shipping_address"`
	BillingAddress  Address     `bson:"billing_address" json:"billing_address"`
	TotalPrice      float64     `bson:"total_price" json:"total_price"`
	Status          OrderStatus `bson:"status" json:"status"`
	CreatedAt       time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `bson:"updated_at" json:"updated_at"`
	//IsActive        bool        `bson:"is_active" json:"is_active"`
}

type OrderItem struct {
	ProductId   string  `bson:"product_id" json:"product_id"`
	ProductName string  `bson:"product_name" json:"product_name"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	UnitPrice   float64 `bson:"unit_price" json:"unit_price"`
}

type Address struct {
	Id      string `bson:"address_id,omitempty" json:"address_id"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	ZipCode string `bson:"zip_code" json:"zip_code"`
}
