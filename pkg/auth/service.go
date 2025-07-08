package auth

import "be-dashboard-nba/pkg/entities"

type Service interface {
	GetUserByEmail(email string) (*entities.User, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) GetUserByEmail(email string) (*entities.User, error) {
	return s.repository.GetUserByEmail(email)
}
