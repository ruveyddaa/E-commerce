package internal

import (
	"tesodev-korpes/OrderService/internal/types"
	"time"

	"github.com/google/uuid"
)

func FromCreateOrderRequest(req *types.CreateOrderRequestModel) *types.Order {
	if req == nil {
		return nil
	}

	items := make([]types.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = types.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	return &types.Order{
		Id:              uuid.NewString(),
		CustomerId:      req.CustomerId,
		Items:           items,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		TotalPrice:      calculateTotalPrice(items),
		Status:          types.OrderOrdered,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		//IsActive:        true,
	}
}

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
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		//IsActive:        order.IsActive,
	}
}

func ToOrderWithCustomerResponse(order *types.OrderResponseModel, customer *types.CustomerResponseModel) *types.OrderWithCustomerResponse {
	if order == nil || customer == nil {
		return nil
	}

	return &types.OrderWithCustomerResponse{
		OrderResponseModel: *order,
		Customer:           *customer,
	}
}
