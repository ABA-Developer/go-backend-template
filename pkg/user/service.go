package user

import "be-dashboard-nba/pkg/entities"

type Service interface {
	CreateUser(user *entities.User) (*entities.User, error)
	ReadUser(id int64) (*entities.User, error)
	UpdateUser(user *entities.User) (*entities.User, error)
	DeleteUser(id int64) error
	GetUsers() ([]entities.User, error)
	GetUserById(id int64) (*entities.User, error)
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

func (s *service) CreateUser(user *entities.User) (*entities.User, error) {
	return s.repository.CreateUser(user)
}

func (s *service) ReadUser(id int64) (*entities.User, error) {
	return s.repository.ReadUser(id)
}

func (s *service) UpdateUser(user *entities.User) (*entities.User, error) {
	return s.repository.UpdateUser(user)
}

func (s *service) DeleteUser(id int64) error {
	return s.repository.DeleteUser(id)
}

func (s *service) GetUsers() ([]entities.User, error) {
	return s.repository.GetAllUser()
}

func (s *service) GetUserById(id int64) (*entities.User, error) {
	return s.repository.GetUserById(id)
}

func (s *service) GetUserByEmail(email string) (*entities.User, error) {
	return s.repository.GetUserByEmail(email)
}
