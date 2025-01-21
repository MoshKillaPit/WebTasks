package services

import (
	"WebTasks/internal/models"
	"WebTasks/internal/repositories"
	"errors"
)

type UserService interface {
	Create(user models.User) (models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id int) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int) error
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

func (s *UserServiceImpl) GetByID(id int) (models.User, error) {
	return s.repo.GetById(id)
}

func (s *UserServiceImpl) Update(user models.User) (models.User, error) {
	return s.repo.Update(user)
}

func (s *UserServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}
