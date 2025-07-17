package internal

import (
	"context"
	"tesodev-korpes/CustomerService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{
		collection: col,
	}
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.Customer, error) {
	var customer *types.Customer
	return customer, nil
}

func (r *Repository) Create(ctx context.Context, customer *types.Customer) (primitive.ObjectID, error) {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	customer.IsActive = true

	result, err := r.collection.InsertOne(ctx, customer)

	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrNilDocument
	}

	return insertedID, nil
}

func (r *Repository) Update(ctx context.Context, id string, update interface{}) error {
	// Placeholder method
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	// Placeholder method
	return nil
}
