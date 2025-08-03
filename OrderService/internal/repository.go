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
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *Repository) Cancel(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     types.OrderCanceled,
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	filter := bson.M{"_id": id}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
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

func (r *Repository) UpdateStatusByID(ctx context.Context, id string, status types.OrderStatus) error {
	var isActive bool
	switch status {
	case types.OrderCanceled, types.OrderDelivered:
		isActive = false
	default:
		isActive = true
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"is_active":  isActive,
			"updated_at": time.Now(),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *Repository) GetAllOrders(ctx context.Context, findOptions *options.FindOptions) ([]types.Order, error) {
	var orders []types.Order

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
