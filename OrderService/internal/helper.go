package internal

import (
	"tesodev-korpes/OrderService/internal/types"
	"time"

	"github.com/google/uuid"
)

func ToOrderResponse(order *types.Order) *types.OrderResponseModel {
	if order == nil {
		return nil
	}

	return &types.OrderResponseModel{
		Id:              order.Id,
		CustomerId:      order.CustomerId,
		Items:           order.Items,
		ShippingAddress: order.ShippingAddress,
		BillingAddress:  order.BillingAddress,
		TotalPrice:      order.TotalPrice,
		Status:          order.Status,
		PaymentStatus:   order.PaymentStatus,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		IsActive:        order.IsActive,
	}
}

func FromCreateOrderRequest(req *types.CreateOrderRequestModel) *types.Order {
	if req == nil {
		return nil
	}

	return &types.Order{
		Id:              uuid.NewString(),
		CustomerId:      req.CustomerId,
		Items:           req.Items,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Status:          types.OrderPending,
		PaymentStatus:   types.PaymentUnpaid,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsActive:        true,
	}
}

func FromUpdateOrderRequest(order *types.Order, req *types.UpdateOrderRequestModel) *types.Order {
	if order == nil || req == nil {
		return order
	}

	if req.Items != nil {
		order.Items = req.Items
	}
	if req.ShippingAddress != nil {
		order.ShippingAddress = *req.ShippingAddress
	}
	if req.BillingAddress != nil {
		order.BillingAddress = *req.BillingAddress
	}
	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.PaymentStatus != nil {
		order.PaymentStatus = *req.PaymentStatus
	}
	if req.IsActive != nil {
		order.IsActive = *req.IsActive
	}

	order.UpdatedAt = time.Now()

	return order
}
