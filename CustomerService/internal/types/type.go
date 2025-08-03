package types

import (
	"time"
)

type CreateCustomerRequestModel struct {
	FirstName string    `bson:"first_name" json:"FirstName" validate:"required,min=2,max=50"`
	LastName  string    `bson:"last_name" json:"last_name" validate:"required,min=2,max=50"`
	Email     string    `bson:"email" json:"email" validate:"required,dive,keys,required,email,endkeys,required"`
	Phone     []Phone   `bson:"phone" json:"phone" validate:"required,dive"`
	Address   []Address `bson:"address" json:"address" validate:"required,dive"`
	Password  string    `bson:"password" json:"password" validate:"required"`
}

type UpdateCustomerRequestModel struct {
	FirstName string    `bson:"first_name,omitempty" json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  string    `bson:"last_name,omitempty" json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
	Email     string    `bson:"email,omitempty" json:"email,omitempty" validate:"omitempty,dive,keys,required,email,endkeys,required"`
	Phone     []Phone   `bson:"phone,omitempty" json:"phone,omitempty" validate:"omitempty,dive"`
	Address   []Address `bson:"address,omitempty" json:"address,omitempty" validate:"omitempty,dive"`
	Password  string    `bson:"password,omitempty" json:"password,omitempty" validate:"omitempty"`
	IsActive  bool      `bson:"is_active,omitempty" json:"is_active,omitempty"`
}

type CustomerResponseModel struct {
	ID        string    `bson:"_id" json:"id"`
	FirstName string    `bson:"first_name" json:"first_name"`
	LastName  string    `bson:"last_name" json:"last_name"`
	Email     string    `bson:"email" json:"email"`
	Phone     []Phone   `bson:"phone" json:"phone"`
	Address   []Address `bson:"address" json:"address"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Pagination struct {
	Page  int
	Limit int
}
type LoginRequestModel struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required"`
}
