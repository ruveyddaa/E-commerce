package internal

import (
	"context"
	"tesodev-korpes/CustomerService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return nil, err
	}

	customer, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	response := types.CustomerResponseModel{
		ID:        customer.Id,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		IsActive:  customer.IsActive,
		Address:   customer.Address,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	return &response, nil
}

func (s *Service) Create(ctx context.Context, req *types.CreateCustomerRequestModel) (string, error) {
	customer := &types.Customer{
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	id, err := s.repo.Create(ctx, customer)
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (s *Service) Update(ctx context.Context, id string, update interface{}) error {
	return s.repo.Update(ctx, id, update)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, params Pagination) ([]types.CustomerResponseModel, error) {
	skip := (params.Page - 1) * params.Limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(params.Limit))

	customers, err := s.repo.Get(ctx, findOptions)
	if err != nil {
		return nil, err
	}
	
	var responses []types.CustomerResponseModel
	for _, c := range customers {
		customerResponse := ToCustomerResponse(&c) 
		responses = append(responses, *customerResponse)
	}

	return responses, nil
}
