package types

import "time"

type Customer struct {
	Id        int               `bson:"_id" json:"id"`
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
	Id      int    `bson:"address_id" json:"address_id"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	ZipCode string `bson:"zip_code" json:"zip_code"`
}

type Phone struct {
	Id          int `bson:"phone_id" json:"phone_id"`
	PhoneNumber int `bson:"phone_number" json:"phone_number"`
}
