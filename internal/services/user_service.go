package services

import (
	"WebTasks/internal/models"
	"WebTasks/internal/repositories"
	"context"
	"errors"
)

type UserService interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, id int) error
}

type userServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) Create(ctx context.Context, user models.User) (models.User, error) {

	if user.Name == "" || user.Key == "" {
		return models.User{}, errors.New("name and key are required")
	}

	createdUser, err := s.repo.Create(ctx, &user)

	if err != nil {
		return models.User{}, err
	}

	return *createdUser, nil
}

func (s *userServiceImpl) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userServiceImpl) GetByID(ctx context.Context, id int) (models.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (s *userServiceImpl) Update(ctx context.Context, user models.User) (models.User, error) {
	if user.Name == "" {
		return models.User{}, errors.New("name is required")
	}

	updatedUser, err := s.repo.Update(ctx, &user)
	if err != nil {
		return models.User{}, err
	}

	return *updatedUser, nil
}

func (s *userServiceImpl) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
