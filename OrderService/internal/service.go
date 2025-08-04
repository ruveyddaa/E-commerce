// File: service.go
package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"tesodev-korpes/pkg/errorPackage"
	"time"

	"tesodev-korpes/OrderService/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *Service) Create(ctx context.Context, order *types.Order) (string, error) {
	if order.CustomerId == "" {
		return "", errors.New("customerId not found ")
	}

	customer, err := s.fetchCustomerByID(order.CustomerId)
	if err != nil || customer == nil {
		return "", fmt.Errorf("customer control unsuccessful")
	}

	order.Id = uuid.NewString()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = types.OrderOrdered
	order.TotalPrice = calculateTotalPrice(order.Items)

	id, err := s.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.OrderWithCustomerResponse, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}

	customer, err := s.fetchCustomerByID(order.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("customer fetch failed: %w", err)
	}

	return &types.OrderWithCustomerResponse{
		OrderResponseModel: *ToOrderResponse(order),
		Customer:           *customer,
	}, nil
}

func (s *Service) ShipOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		return err
	}

	if order.Status != types.OrderOrdered {
		return errorPackage.InvalidOrderStateWithStatus("ship", string(order.Status))
	}

	err = s.repo.UpdateStatusByID(ctx, id, types.OrderShipped)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeliverOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		return err
	}

	if order.Status != types.OrderShipped {
		return errorPackage.InvalidOrderStateWithStatus("deliver", string(order.Status))
	}

	err = s.repo.UpdateStatusByID(ctx, id, types.OrderDelivered)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		return err
	}

	switch order.Status {
	case types.OrderShipped, types.OrderDelivered, types.OrderCanceled:
		return errorPackage.InvalidOrderStateWithStatus("CANCEL", string(order.Status))

	}

	err = s.repo.UpdateStatusByID(ctx, id, types.OrderCanceled)
	if err != nil {
		return err
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

func (s *Service) fetchCustomerByID(customerID string) (*types.CustomerResponseModel, error) {
	if customerID == "" {
		return nil, errors.New("customerID bos")
	}

	url := fmt.Sprintf("%s/customer/%s", s.customerServiceURL, customerID)
	fmt.Printf("Fetching customer from URL: %s\n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("Customer service response status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Customer service returned non-OK status: %d\n", resp.StatusCode)
		return nil, nil
	}

	var customer types.CustomerResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		fmt.Printf("Error decoding customer response: %v\n", err)
		return nil, err
	}

	fmt.Printf("Successfully fetched customer: %s\n", customer.Id)
	return &customer, nil
}

func calculateTotalPrice(items []types.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}
