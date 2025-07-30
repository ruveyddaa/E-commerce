package types

import (
	"time"
)

type Customer struct {
	Id        string            `bson:"_id,omitempty" json:"id"`
	Password  []byte            `bson:"password" json:"password"`
	FirstName string            `bson:"first_name" json:"first_name"`
	LastName  string            `bson:"last_name" json:"last_name"`
	Email     map[string]string `bson:"email" json:"email"`
	Phone     []Phone           `bson:"phone" json:"phone"`
	Address   []Address         `bson:"address" json:"address"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
	IsActive  bool              `bson:"is_active" json:"is_active"`
}

type Address struct {
	Id      string `bson:"address_id,omitempty" json:"address_id"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	ZipCode string `bson:"zip_code" json:"zip_code"`
}

type Phone struct {
	Id          string `bson:"phone_id,omitempty" json:"phone_id"`
	PhoneNumber int    `bson:"phone_number" json:"phone_number"`
}
