package internal

import (
	"context"
	"errors"
	"fmt"
	"tesodev-korpes/CustomerService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *Service) GetByID(ctx context.Context, id string) (*types.CustomerResponseModel, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, NewBadRequest(fmt.Sprintf("invalid ID format: %s", id))
	}

	customer, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, NewNotFound(fmt.Sprintf("customer not found for ID: %s", id))
		}
		return nil, NewInternal(err.Error())
	}

	return ToCustomerResponse(customer), nil
}

func (s *Service) Create(ctx context.Context, req *types.CreateCustomerRequestModel) (string, error) {
	customer := FromCreateCustomerRequest(req)
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	customer.IsActive = true // todo bir defaoltu atılacak
	id, err := s.repo.Create(ctx, customer)
	if err != nil {
		return "", NewNotFound(fmt.Sprintf("failed to create customer: %v", err))
	} // todo ezgi

	return id.Hex(), nil
}

func (s *Service) Update(ctx context.Context, id string, req *types.UpdateCustomerRequestModel) (*types.Customer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		//return nil, (fmt.Sprintf("invalid ID format: %w", err))
		return nil, NewBadRequest(fmt.Sprintf("invalid ID format: %s", id))
	}
	customer, err := s.repo.GetByID(ctx, objectID)
	// isExist uluştur, dbden glen müşeriyi alma
	if err != nil {
		//return nil, fmt.Errorf("customer not found: %w", err)
		return nil, NewNotFound(fmt.Sprintf("customer not found for ID: %s", id))
	}
	updatedCustomer := FromUpdateCustomerRequest(customer, req)
	// todo tek req ten ilerlet

	err = s.repo.Update(ctx, objectID, updatedCustomer)
	if err != nil {
		return nil, NewInternal("failed to update customer")
	}
	return updatedCustomer, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewBadRequest(fmt.Sprintf("Invalid id format: %s", id))
	}

	_, err = s.repo.GetByID(ctx, objectID)
	if err != nil {
		return NewNotFound(fmt.Sprintf("Customer not found with id %s", id))
	} // todo gereksiz silinecek

	if err := s.repo.Delete(ctx, objectID); err != nil {
		return NewInternal(fmt.Sprintf("Failed to delete customer with id %s", id))
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
		if err == mongo.ErrNoDocuments {
			return nil, NewNotFound("customer not found")
		}
		return nil, NewInternal(err.Error())
	}

	var responses []types.CustomerResponseModel
	for _, c := range customers {
		customerResponse := ToCustomerResponse(&c)
		responses = append(responses, *customerResponse)
	}

	return responses, nil
}
