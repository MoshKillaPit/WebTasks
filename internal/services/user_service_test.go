package services

import (
	"WebTasks/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// MockUserRepository реализует методы UserRepository для тестов.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

// Исправлено: GetByID вместо GetById
func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestUserService_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	validUser := models.User{
		Name: "Valid User",
		Key:  "valid-key",
	}
	invalidUser := models.User{
		Name: "",
		Key:  "",
	}

	// Успешное создание
	mockRepo.On("Create", ctx, &validUser).Return(&validUser, nil)

	result, err := service.Create(ctx, validUser)
	require.NoError(t, err)
	require.Equal(t, validUser, result)
	mockRepo.AssertExpectations(t)

	// Ошибка: пустые имя и ключ
	_, err = service.Create(ctx, invalidUser)
	require.Error(t, err)
	require.Equal(t, "name and key are required", err.Error())
	mockRepo.AssertNotCalled(t, "Create", ctx, &invalidUser)
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	validUser := models.User{
		ID:   1,
		Name: "Valid User",
		Key:  "valid-key",
	}

	// Успешное получение пользователя
	mockRepo.On("GetByID", ctx, validUser.ID).Return(&validUser, nil)

	result, err := service.GetByID(ctx, validUser.ID)
	require.NoError(t, err)
	require.Equal(t, validUser, result)
	mockRepo.AssertExpectations(t)

	// Ошибка: пользователь не найден
	mockRepo.On("GetByID", ctx, 999).Return(nil, errors.New("user not found"))

	_, err = service.GetByID(ctx, 999)
	require.Error(t, err)
	require.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	updatedUser := models.User{
		ID:   1,
		Name: "Updated User",
		Key:  "existing-key",
	}

	// Успешное обновление
	mockRepo.On("Update", ctx, &updatedUser).Return(&updatedUser, nil)

	result, err := service.Update(ctx, updatedUser)
	require.NoError(t, err)
	require.Equal(t, updatedUser, result)
	mockRepo.AssertExpectations(t)

	// Ошибка: пустое имя
	invalidUser := models.User{
		ID:   1,
		Name: "",
	}
	_, err = service.Update(ctx, invalidUser)
	require.Error(t, err)
	require.Equal(t, "name is required", err.Error())
	mockRepo.AssertNotCalled(t, "Update", ctx, &invalidUser)
}

func TestUserService_Delete(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	validUserID := 1

	// Успешное удаление
	mockRepo.On("Delete", ctx, validUserID).Return(nil)

	err := service.Delete(ctx, validUserID)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Ошибка: пользователь не найден
	mockRepo.On("Delete", ctx, 999).Return(errors.New("user not found"))

	err = service.Delete(ctx, 999)
	require.Error(t, err)
	require.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
