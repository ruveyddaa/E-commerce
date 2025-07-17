package internal

import (
	"context"
	"tesodev-korpes/CustomerService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	Page  int
	Limit int
}

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{
		collection: col,
	}
}

func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*types.Customer, error) {
	var customer types.Customer
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
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

func (r *Repository) Get(ctx context.Context, opt *options.FindOptions) ([]types.Customer, error) {
	var customers []types.Customer

	cursor, err := r.collection.Find(ctx, bson.M{}, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}
