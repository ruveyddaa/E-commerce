package internal

import (
	"tesodev-korpes/CustomerService/internal/types"
)

// here is an example what this helper method that does data casting from db model to response model
// the return statement that I commented out repreents an introduction that how you can implement it
// you can delete after you'd completed the helper method, its a placeholder put here just to prevent getting errors at
// the beginning
// ToCustomerResponse converts Customer entity to CustomerResponseModel
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

// ToCustomer converts CustomerResponseModel back to Customer entity
func ToCustomer(resp *types.CustomerResponseModel) *types.Customer {
	if resp == nil {
		return nil
	}
	return &types.Customer{
		Id:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		Phone:     resp.Phone,
		Address:   resp.Address,
		IsActive:  resp.IsActive,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

// FromCreateCustomerRequest converts CreateCustomerRequestModel to Customer entity
func FromCreateCustomerRequest(req *types.CreateCustomerRequestModel) *types.Customer {
	if req == nil {
		return nil
	}
	return &types.Customer{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Password:  req.Password,
		IsActive:  req.IsActive,
	}
}

// FromUpdateCustomerRequest updates existing Customer entity with UpdateCustomerRequestModel fields
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
	if req.Email != nil {
		customer.Email = req.Email
	}
	if req.Phone != nil {
		customer.Phone = req.Phone
	}
	if req.Address != nil {
		customer.Address = req.Address
	}
	if req.Password != nil {
		customer.Password = req.Password
	}
	if !req.IsActive {
		customer.IsActive = req.IsActive
	}
	return customer
}
