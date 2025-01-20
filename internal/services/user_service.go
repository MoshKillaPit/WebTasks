package services

import (
	"WebTasks/internal/models"
	"WebTasks/internal/repositories"
	"github.com/pkg/errors"
)

type UserService interface {
	Create(user models.User) (models.User, error)
	GetAll() ([]models.User, error)
}

type UserServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Create(user models.User) (models.User, error) {
	if user.Name == "" || user.Key == "" {
		return models.User{}, errors.New("name and key are required")
	}
	return s.repo.Create(user)
}

func (s *UserServiceImpl) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}
