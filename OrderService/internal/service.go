package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"tesodev-korpes/OrderService/internal/types"
	// "time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
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

func (s *Service) DeleteOrderByID(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	return s.repo.Delete(id)
}
func (s *Service) Create(ctx context.Context, order *types.Order) (string, error) {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.IsActive = true

	return s.repo.Create(ctx, order)
}
