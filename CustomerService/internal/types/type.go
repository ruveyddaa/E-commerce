package types

import (
	"time"
)

type CreateCustomerRequestModel struct {
	FirstName    string       `json:"first_name" validate:"required,min=2,max=50"`
	LastName     string       `json:"last_name" validate:"required,min=2,max=50"`
	Email        string       `json:"email" validate:"required,email"`
	Phone        []Phone      `json:"phone" validate:"required,dive"`
	Address      []Address    `json:"address" validate:"required,dive"`
	Password     string       `json:"password" validate:"required"`
	System       System       `json:"system" validate:"required,dive"`
	Subscription Subscription `json:"subscription" validate:"required,dive"`
}

type UpdateCustomerRequestModel struct {
	FirstName    string       `json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
	LastName     string       `json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
	Email        string       `json:"email" validate:"required,email"`
	Phone        []Phone      `json:"phone,omitempty" validate:"omitempty,dive"`
	Address      []Address    `json:"address,omitempty" validate:"omitempty,dive"`
	Password     string       `json:"password,omitempty" validate:"omitempty"`
	IsActive     bool         `json:"is_active,omitempty"`
	System       System       `json:"system,omitempty" validate:"omitempty,dive"`
	Subscription Subscription `json:"subscription,omitempty" validate:"omitempty,dive"`
}

type CustomerResponseModel struct {
	ID           string       `json:"id"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	Email        string       `json:"email"`
	Phone        []Phone      `json:"phone"`
	Address      []Address    `json:"address"`
	IsActive     bool         `json:"is_active"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	System       System       `json:"system"`
	Subscription Subscription `json:"subscription"`
}

type Pagination struct {
	Page  int
	Limit int
}

type LoginRequestModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token   string                 `json:"token"`
	User    *CustomerResponseModel `json:"user"`
	Message string                 `json:"message"`
}

type VerifiedUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type VerifyTokenResponse struct {
	Message string       `json:"message"`
	User    VerifiedUser `json:"user"`
}
