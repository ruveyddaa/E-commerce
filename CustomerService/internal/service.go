package internal

import (
	"context"
	"errors"
	"fmt"
	"tesodev-korpes/CustomerService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

/*
	func (s *Service) GetByID(ctx context.Context, id string) (*types.CustomerResponseModel, error) {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		customer, err := s.repo.GetByID(ctx, objectID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf("customer not found for ID: %s", id)
			}
			return nil, fmt.Errorf("failed to get customer: %w", err)
		}

		return ToCustomerResponse(customer), nil
	}
*/
func (s *Service) GetByID(ctx context.Context, id string) (*types.CustomerResponseModel, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("customer not found for ID: %s", id)
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return ToCustomerResponse(customer), nil
}

/*
	func (s *Service) Create(ctx context.Context, req *types.CreateCustomerRequestModel) (string, error) {
		customer := FromCreateCustomerRequest(req)
		customer.CreatedAt = time.Now()
		customer.UpdatedAt = time.Now()

		id, err := s.repo.Create(ctx, customer)

		if err != nil {
			return "", fmt.Errorf("failed to create customer: %w", err)

		}

		return id.Hex(), nil
	}
*/
func (s *Service) Create(ctx context.Context, req *types.CreateCustomerRequestModel) (string, error) {
	customer := FromCreateCustomerRequest(req)
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	id, err := s.repo.Create(ctx, customer)
	if err != nil {
		return "", fmt.Errorf("failed to create customer: %w", err)
	}
	return id, nil
}

/*
	func (s *Service) Update(ctx context.Context, id string, req *types.UpdateCustomerRequestModel) (*types.Customer, error) {
		objectID, err := uuid.UUID(id)
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %s", id)
		}
		customer, err := s.repo.GetByID(ctx, objectID)
		// isExist uluştur, dbden glen müşeriyi alma
		if err != nil {
			return nil, fmt.Errorf("customer not found for ID: %s", id)
		}
		updatedCustomer := FromUpdateCustomerRequest(customer, req)
		// todo tek req ten ilerlet

		err = s.repo.Update(ctx, objectID, updatedCustomer)
		if err != nil {
			return nil, errors.New("failed to update customer")
		}
		return updatedCustomer, nil
	}
*/
func (s *Service) Update(ctx context.Context, id string, req *types.UpdateCustomerRequestModel) (*types.Customer, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("customer not found for ID: %s", id)
	}

	updatedCustomer := FromUpdateCustomerRequest(customer, req)

	err = s.repo.Update(ctx, id, updatedCustomer)
	if err != nil {
		return nil, errors.New("failed to update customer")
	}

	return updatedCustomer, nil
}

/*
	func (s *Service) Delete(ctx context.Context, id string) error {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("Invalid id format: %s", id)
		}

		_, err = s.repo.GetByID(ctx, objectID)
		if err != nil {
			return fmt.Errorf("Customer not found with id %s", id)
		} // todo gereksiz silinecek

		if err := s.repo.Delete(ctx, objectID); err != nil {
			return fmt.Errorf("Failed to delete customer with id %s", id)
		}
		return nil
	}
*/
func (s *Service) Delete(ctx context.Context, id string) error {
	// 1. Müşteri var mı kontrol et
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("customer not found with id %s", id)
	}

	// 2. Silme işlemi
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete customer with id %s", id)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, params types.Pagination) ([]types.CustomerResponseModel, error) {
	skip := (params.Page - 1) * params.Limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(params.Limit))

	customers, err := s.repo.Get(ctx, findOptions)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	var responses []types.CustomerResponseModel
	for _, c := range customers {
		customerResponse := ToCustomerResponse(&c)
		responses = append(responses, *customerResponse)
	}

	return responses, nil
}
