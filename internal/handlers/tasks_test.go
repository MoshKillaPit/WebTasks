package handlers_test

import (
	"WebTasks/internal/handlers"
	"WebTasks/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskService - мок для интерфейса TaskService
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) Create(ctx context.Context, task models.Task) (models.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskService) GetByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) GetAll(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskService) Update(ctx context.Context, task models.Task) (models.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --------------------------------------------------------------------------------------
// ТЕСТЫ НА УСПЕШНОЕ ПОВЕДЕНИЕ
// --------------------------------------------------------------------------------------

func TestHandler_GetTasks(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	expectedTasks := []models.Task{
		{ID: 1, Name: "Task 1", Status: "Pending"},
		{ID: 2, Name: "Task 2", Status: "Completed"},
	}

	mockService.On("GetAll", mock.Anything).Return(expectedTasks, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	handler.GetTasks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualTasks []models.Task
	err := json.NewDecoder(rr.Body).Decode(&actualTasks)
	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, actualTasks)

	mockService.AssertExpectations(t)
}

func TestHandler_GetTaskByID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	expectedTask := &models.Task{ID: 1, Name: "Task 1", Status: "Pending"}

	mockService.On("GetByID", mock.Anything, 1).Return(expectedTask, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	// симуляция того, что {id} = "1"
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.GetTaskByID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualTask models.Task
	err := json.NewDecoder(rr.Body).Decode(&actualTask)
	assert.NoError(t, err)
	assert.Equal(t, *expectedTask, actualTask)

	mockService.AssertExpectations(t)
}

func TestHandler_CreateTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	inputTask := models.Task{Name: "New Task", Status: "Pending"}
	expectedTask := models.Task{ID: 1, Name: "New Task", Status: "Pending"}

	mockService.On("Create", mock.Anything, inputTask).Return(expectedTask, nil)

	body, _ := json.Marshal(inputTask)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.CreateTask(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var actualTask models.Task
	err := json.NewDecoder(rr.Body).Decode(&actualTask)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, actualTask)

	mockService.AssertExpectations(t)
}

func TestHandler_UpdateTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	inputTask := models.Task{ID: 1, Name: "Updated Task", Status: "Completed"}
	expectedTask := models.Task{ID: 1, Name: "Updated Task", Status: "Completed"}

	mockService.On("Update", mock.Anything, inputTask).Return(expectedTask, nil)

	body, _ := json.Marshal(inputTask)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.UpdateTask(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualTask models.Task
	err := json.NewDecoder(rr.Body).Decode(&actualTask)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, actualTask)

	mockService.AssertExpectations(t)
}

func TestHandler_DeleteTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	mockService.On("Delete", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.DeleteTask(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())

	mockService.AssertExpectations(t)
}

// --------------------------------------------------------------------------------------
// ТЕСТЫ НА ОШИБОЧНОЕ ПОВЕДЕНИЕ (ПОЛУЧЕНИЕ 100% ПОКРЫТИЯ)
// --------------------------------------------------------------------------------------

func TestHandler_GetTasks_ServiceError(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	mockService.On("GetAll", mock.Anything).Return([]models.Task{}, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	handler.GetTasks(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to fetch tasks")

	mockService.AssertExpectations(t)
}

func TestHandler_GetTaskByID_InvalidID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	// Передаём некорректный ID, например "abc"
	req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	handler.GetTaskByID(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid task ID")

	mockService.AssertNotCalled(t, "GetByID")
}

func TestHandler_GetTaskByID_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	mockService.On("GetByID", mock.Anything, 999).Return((*models.Task)(nil), errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()

	handler.GetTaskByID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Task not found")

	mockService.AssertExpectations(t)
}

func TestHandler_CreateTask_InvalidJSON(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	// Намеренно некорректный JSON (пропущена кавычка или скобка)
	body := []byte(`{"Name":"Bad Task", "Status":`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.CreateTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")

	// Сервис вызываться не должен
	mockService.AssertNotCalled(t, "Create")
}

func TestHandler_CreateTask_ServiceError(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	inputTask := models.Task{Name: "New Task", Status: "Pending"}

	mockService.On("Create", mock.Anything, inputTask).Return(models.Task{}, errors.New("service create error"))

	body, _ := json.Marshal(inputTask)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.CreateTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "service create error")

	mockService.AssertExpectations(t)
}

func TestHandler_UpdateTask_InvalidID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	body, _ := json.Marshal(models.Task{Name: "Should fail"})
	req := httptest.NewRequest(http.MethodPut, "/tasks/abc", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	handler.UpdateTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid task ID")

	mockService.AssertNotCalled(t, "Update")
}

func TestHandler_UpdateTask_InvalidJSON(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer([]byte(`{"Name":"Bad JSON"`)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.UpdateTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")

	mockService.AssertNotCalled(t, "Update")
}

func TestHandler_UpdateTask_ServiceError(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	inputTask := models.Task{ID: 1, Name: "Error Task", Status: "New"}

	mockService.On("Update", mock.Anything, inputTask).Return(models.Task{}, errors.New("update error"))

	body, _ := json.Marshal(inputTask)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.UpdateTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "update error")

	mockService.AssertExpectations(t)
}

func TestHandler_DeleteTask_InvalidID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	handler.DeleteTask(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid task ID")

	mockService.AssertNotCalled(t, "Delete")
}

func TestHandler_DeleteTask_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := handlers.NewHandler(mockService)

	mockService.On("Delete", mock.Anything, 999).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr := httptest.NewRecorder()

	handler.DeleteTask(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Task not found")

	mockService.AssertExpectations(t)
}
