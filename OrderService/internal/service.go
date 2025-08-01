package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"tesodev-korpes/pkg"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"tesodev-korpes/OrderService/internal/types"
	// "time"
)

type Service struct {
	repo               *Repository
	httpClient         *http.Client
	customerServiceURL string
}

func NewService(repo *Repository, customerServiceURL string) *Service {
	return &Service{
		repo:               repo,
		httpClient:         &http.Client{Timeout: 5 * time.Second},
		customerServiceURL: customerServiceURL,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.OrderResponseModel, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("order not found for ID: %s", id)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return ToOrderResponse(order), nil
}

func (s *Service) Create(ctx context.Context, order *types.Order) (string, error) {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	//order.IsActive = true

	return s.repo.Create(ctx, order)
}

func (s *Service) ShipOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("order not found for ID: %s", id)
		}
		return err
	}

	if order.Status != types.OrderOrdered {
		return pkg.InvalidOrderStateWithStatus("ship", string(order.Status))
	}

	err = s.repo.UpdateStatusByID(ctx, id, types.OrderShipped)
	if err != nil {
		return fmt.Errorf("failed to update order status to SHIPPED: %w", err)
	}

	return nil
}

func (s *Service) DeliverOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("order not found for ID: %s", id)
		}
		return err
	}

	if order.Status != types.OrderShipped {
		return pkg.InvalidOrderStateWithStatus("deliver", string(order.Status))
	}

	err = s.repo.UpdateStatusByID(ctx, id, types.OrderDelivered)
	if err != nil {
		return fmt.Errorf("failed to update order status to DELIVERED: %w", err)
	}

	return nil
}

func calculateTotalPrice(items []types.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}

func (s *Service) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("order not found for ID: %s", id)
		}
		return err
	}

	switch order.Status {
	case types.OrderShipped, types.OrderDelivered, types.OrderCanceled:
		return pkg.InvalidOrderStateWithStatus("CANCEL", string(order.Status))

	}

	err = s.repo.UpdateStatusByID(ctx, id, types.OrderCanceled)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

func (s *Service) GetAllOrders(ctx context.Context, pagination types.Pagination) ([]*types.OrderResponseModel, error) {
	skip := (pagination.Page - 1) * pagination.Limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.Limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	orders, err := s.repo.GetAllOrders(ctx, findOptions)
	if err != nil {
		return nil, err
	}

	var response []*types.OrderResponseModel
	for _, order := range orders {
		resp := ToOrderResponse(&order)
		if resp != nil {
			response = append(response, resp)
		}
	}

	return response, nil

}
