package types

import (
	"time"
)

type Customer struct {
	Id        string    `bson:"_id,omitempty"`
	Password  string    `bson:"password"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Email     string    `bson:"email"`
	Phone     []Phone   `bson:"phone"`
	Address   []Address `bson:"address"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	IsActive  bool      `bson:"is_active"`
	Role      string    `bson:"role"`
	Token     string    `bson:"token"`
}

type Address struct {
	Id      string `bson:"address_id,omitempty"`
	City    string `bson:"city"`
	State   string `bson:"state"`
	ZipCode string `bson:"zip_code"`
}

type Phone struct {
	Id          string `bson:"phone_id,omitempty"`
	PhoneNumber int    `bson:"phone_number"`
}
