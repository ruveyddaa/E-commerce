package types

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderShipped   OrderStatus = "SHIPPED"
	OrderDelivered OrderStatus = "DELIVERED"
	OrderCanceled  OrderStatus = "CANCELED"
	OrderOrdered   OrderStatus = "ORDERED"
)
