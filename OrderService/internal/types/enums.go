package types

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderShipped   OrderStatus = "SHIPPED"
	OrderDelivered OrderStatus = "DELIVERED"
	OrderCanceled  OrderStatus = "CANCELED"
)

type PaymentStatus string

const (
	PaymentUnpaid   PaymentStatus = "UNPAID"
	PaymentPaid     PaymentStatus = "PAID"
	PaymentRefunded PaymentStatus = "REFUNDED"
)
