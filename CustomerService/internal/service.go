package internal

import (
	"context"
	"errors"
	"fmt"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/auth"
	"tesodev-korpes/pkg/errorPackage"
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
func (s *Service) Login(ctx context.Context, email, password, correlationID string) (string, *types.Customer, error) {
	customer, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil, errorPackage.NewNotFound("404101")
		}
		pkg.LogErrorWithCorrelation(err, correlationID)
		return "", nil, errorPackage.NewInternal("500101", err)
	}

	valid, err := auth.VerifyPassword(password, customer.Password)
	if err != nil {
		pkg.LogErrorWithCorrelation(err, correlationID)
		return "", nil, errorPackage.NewInternal("500101", err)
	}
	if !valid {
		return "", nil, errorPackage.NewUnauthorized("404201")
	}

	token, err := auth.GenerateJWT(customer.Id)
	if err != nil {
		pkg.LogErrorWithCorrelation(err, correlationID)
		return "", nil, errorPackage.NewInternal("500101", err)
	}

	return token, customer, nil
}
func (s *Service) GetByEmail(ctx context.Context, email string) (*types.Customer, error) {
	customer, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err // yukarıda 404 olarak dönecek
		}
		return nil, err
	}

	return customer, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.CustomerResponseModel, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return ToCustomerResponse(customer), nil
}

func (s *Service) Create(ctx context.Context, req *types.CreateCustomerRequestModel) (string, error) {
	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	req.Password = string(hashedPwd)

	customer := FromCreateCustomerRequest(req)
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	customer.Role = "non-premium" // role atması yapılıyor

	id, err := s.repo.Create(ctx, customer)
	if err != nil {
		return "", fmt.Errorf("failed to create customer: %w", err)
	}

	return id, nil
}

func (s *Service) Update(ctx context.Context, id string, customer *types.Customer) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("customer not found for ID: %s", id)
	}
	return s.repo.Update(ctx, id, customer)
}
func (s *Service) Delete(ctx context.Context, id string) error {

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
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
		return nil, err
	}

	var responses []types.CustomerResponseModel
	for _, c := range customers {
		customerResponse := ToCustomerResponse(&c)
		responses = append(responses, *customerResponse)
	}

	return responses, nil
}
