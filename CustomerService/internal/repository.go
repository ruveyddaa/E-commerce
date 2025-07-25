package internal

import (
	"context"
	"fmt"
	"tesodev-korpes/CustomerService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, customer *types.Customer) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{ // todo mongo operatorlerle çözebilrisin
			"first_name": customer.FirstName,
			"last_name":  customer.LastName,
			"email":      customer.Email,
			"phone":      customer.Phone,
			"address":    customer.Address,
			"password":   customer.Password,
			"is_active":  customer.IsActive,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, opt *options.FindOptions) ([]types.Customer, error) {
	var customers []types.Customer

	cursor, err := r.collection.Find(ctx, bson.M{}, opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document not found")
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}
