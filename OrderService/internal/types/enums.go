package types

type OrderStatus string

const (
	OrderOrdered   OrderStatus = "ORDERED"
	OrderShipped   OrderStatus = "SHIPPED"
	OrderDelivered OrderStatus = "DELIVERED"
	OrderCanceled  OrderStatus = "CANCELED"
)
