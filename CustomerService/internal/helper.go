package internal

import (
	"tesodev-korpes/CustomerService/internal/types"

	"github.com/google/uuid"
)

// here is an example what this helper method that does data casting from db model to response model
// the return statement that I commented out repreents an introduction that how you can implement it
// you can delete after you'd completed the helper method, its a placeholder put here just to prevent getting errors at
// the beginning

func ToCustomerResponse(customer *types.Customer) *types.CustomerResponseModel {
	if customer == nil {
		return nil
	}
	return &types.CustomerResponseModel{
		ID:        customer.Id,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Address:   customer.Address,
		IsActive:  customer.IsActive,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}
func FromCreateCustomerRequest(req *types.CreateCustomerRequestModel) *types.Customer {
	if req == nil {
		return nil
	}

	addresses := make([]types.Address, len(req.Address))
	for i, addr := range req.Address {
		addresses[i] = types.Address{
			Id:      uuid.New().String(),
			City:    addr.City,
			State:   addr.State,
			ZipCode: addr.ZipCode,
		}
	}

	phones := make([]types.Phone, len(req.Phone))
	for i, ph := range req.Phone {
		phones[i] = types.Phone{
			Id:          uuid.New().String(),
			PhoneNumber: ph.PhoneNumber,
		}
	}

	return &types.Customer{
		Id:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     phones,
		Address:   addresses,
		Password:  req.Password,
		IsActive:  true,
	}
}
func FromUpdateCustomerRequest(customer *types.Customer, req *types.UpdateCustomerRequestModel) *types.Customer {
	if customer == nil || req == nil {
		return customer
	}
	if req.FirstName != "" {
		customer.FirstName = req.FirstName
	}
	if req.LastName != "" {
		customer.LastName = req.LastName
	}

	if req.Phone != nil {
		customer.Phone = req.Phone
	}
	if req.Address != nil {
		customer.Address = req.Address
	}
	if !req.IsActive {
		customer.IsActive = req.IsActive
	}
	return customer
}
