package types

import "time"

type CreateCustomerRequestModel struct {
	Password  []byte            `bson:"password" json:"password"`
	FirstName string            `bson:"first_name" json:"first_name" validate:"required,min=2,max=50"`
	LastName  string            `bson:"last_name" json:"last_name" validate:"required,min=2,max=50"`
	Email     map[string]string `bson:"email" json:"email" validate:"required,dive,keys,required,email,endkeys,required"`
	Phone     []Phone           `bson:"phone" json:"phone" validate:"required,dive"`
	Address   []Address         `bson:"address" json:"address" validate:"required,dive"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at" validate:"omitempty"`
	IsActive  bool              `bson:"is_active" json:"is_active" validate:"omitempty"`
}

type UpdateCustomerRequestModel struct {
	Password  []byte            `bson:"password" json:"password"`
	FirstName string            `bson:"first_name" json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string            `bson:"last_name" json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     map[string]string `bson:"email" json:"email" validate:"omitempty,dive,keys,required,email,endkeys,required"`
	Phone     []Phone           `bson:"phone" json:"phone" validate:"omitempty,dive"`
	Address   []Address         `bson:"address" json:"address" validate:"omitempty,dive"`
	IsActive  bool              `bson:"is_active" json:"is_active" validate:"omitempty"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at" validate:"omitempty"`
}

type CustomerResponseModel struct {
	Password  []byte            `bson:"password" json:"password"`
	FirstName string            `bson:"first_name" json:"first_name" `
	LastName  string            `bson:"last_name" json:"last_name" `
	Email     map[string]string `bson:"email" json:"email" `
	Phone     []Phone           `bson:"phone" json:"phone"`
	Address   []Address         `bson:"address" json:"address"`
	UpdatedAt time.Time         `bson:"created_at" json:"created_at"`
	IsActive  bool              `bson:"is_active" json:"is_active"`
}
