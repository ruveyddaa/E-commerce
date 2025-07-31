package internal

import (
	"context"
	"errors"
	"fmt"
	"tesodev-korpes/OrderService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}
func (r *Repository) GetByID(ctx context.Context, id string) (*types.Order, error) {
	// 1. String ID'yi MongoDB ObjectID'ye dönüştür
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("geçersiz id formatı: %w", err)
	}

	// 2. Filter ile eşleşen dökümanı bul ve decode et
	var order types.Order
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	// 3. Order bulunduysa geri döndür
	return &order, nil
}

func (r *Repository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}
func (r *Repository) Create(ctx context.Context, order *types.Order) (string, error) {
	res, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	// ObjectID'yi string olarak döndür
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}
