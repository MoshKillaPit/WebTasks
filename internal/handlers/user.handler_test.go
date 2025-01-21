package handlers_test

import (
	"WebTasks/internal/handlers"
	"WebTasks/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService - мок для интерфейса UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id int) (models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_GetUsers(t *testing.T) {
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	expectedUsers := []models.User{
		{ID: 1, Name: "Alice", Key: "key123"},
		{ID: 2, Name: "Bob", Key: "key456"},
	}

	mockService.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualUsers []models.User
	err := json.NewDecoder(rr.Body).Decode(&actualUsers)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, actualUsers)

	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	inputUser := models.User{Name: "Alice", Key: "key123"}
	expectedUser := models.User{ID: 1, Name: "Alice", Key: "key123"}

	mockService.On("Create", mock.Anything, inputUser).Return(expectedUser, nil)

	body, _ := json.Marshal(inputUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var actualUser models.User
	err := json.NewDecoder(rr.Body).Decode(&actualUser)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)

	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidBody(t *testing.T) {
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid body")))
	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUsers_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	// Возвращаем пустой слайс пользователей и ошибку
	mockService.On("GetAll", mock.Anything).Return([]models.User{}, context.DeadlineExceeded)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler.GetUsers(rr, req)

	// Ожидаем, что обработка завершится ошибкой
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to fetch users")

	mockService.AssertExpectations(t)
}
