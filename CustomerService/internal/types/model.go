package types

import "time"

type Customer struct {
	Id        string    `bson:"_id" json:"id"`
	FirstName string    `bson:"first_name" json:"first_name"`
	LastName  string    `bson:"last_name" json:"last_name"`
	Email     string    `bson:"email" json:"email"`
	Phone     string    `bson:"phone" json:"phone"`
	Address   string    `bson:"address" json:"address"`
	City      string    `bson:"city" json:"city"`
	State     string    `bson:"state" json:"state"`
	ZipCode   string    `bson:"zip_code" json:"zip_code"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
