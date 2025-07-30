package internal

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) DeleteOrderByID(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	return s.repo.Delete(id)
}
