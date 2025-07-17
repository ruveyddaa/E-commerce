package internal

import (
	"context"
	"fmt"
	"tesodev-korpes/CustomerService/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceCustomerRequestModel struct {
	FirstName string
	LastName  string
	Email     map[string]string
	Phone     []types.Phone
	Address   []types.Address
	Password  []byte
	IsActive  bool
}

type ServiceCustomerResponseModel struct {
	ID        primitive.ObjectID
	FirstName string
	LastName  string
	Email     map[string]string
	Phone     []types.Phone
	Address   []types.Address
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewBadRequest(fmt.Sprintf("Invalid id format: %s", id))
	}

	_, err = s.repo.GetByID(ctx, objectID)
	if err != nil {
		return NewNotFound(fmt.Sprintf("Customer not found with id %s", id))
	}

	if err := s.repo.DeleteByObjectID(ctx, objectID); err != nil {
		return NewInternal(fmt.Sprintf("Failed to delete customer with id %s", id))
	}
	return nil
}

func (s *Service) Get(ctx context.Context, params Pagination) ([]ServiceCustomerResponseModel, error) {
	skip := (params.Page - 1) * params.Limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(params.Limit))

	customers, err := s.repo.Get(ctx, findOptions)
	if err != nil {
		return nil, err
	}

	var responses []ServiceCustomerResponseModel
	for _, c := range customers {
		responses = append(responses, ServiceCustomerResponseModel{
			ID:        c.Id,
			FirstName: c.FirstName,
			LastName:  c.LastName,
			Email:     c.Email,
			Phone:     c.Phone,
			Address:   c.Address,
			IsActive:  c.IsActive,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	return responses, nil
}
