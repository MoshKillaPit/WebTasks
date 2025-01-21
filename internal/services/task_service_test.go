package services

import (
	"WebTasks/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// MockTaskRepository - мок для TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	args := m.Called(ctx, task)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) GetById(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	args := m.Called(ctx, task)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Тест для Create
func TestTaskService_Create(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	ctx := context.Background()
	validTask := models.Task{
		Name:   "Test Task",
		Status: "Pending",
		Time:   time.Now(),
		Due:    time.Now().Add(24 * time.Hour),
	}
	invalidTask := models.Task{
		Name: "",
	}

	// Успешное создание
	mockRepo.On("Create", ctx, &validTask).Return(&validTask, nil)

	result, err := service.Create(ctx, validTask)

	require.NoError(t, err)
	require.Equal(t, validTask, result)
	mockRepo.AssertExpectations(t)

	// Ошибка: пустое имя
	_, err = service.Create(ctx, invalidTask)
	require.Error(t, err)
	require.Equal(t, "task name is required", err.Error())
	mockRepo.AssertNotCalled(t, "Create", ctx, &invalidTask)

	// Ошибка: слишком длинное имя
	invalidTask = models.Task{Name: "This is a very long task name that exceeds fifty characters"}
	_, err = service.Create(ctx, invalidTask)
	require.Error(t, err)
	require.Equal(t, "task name is too long", err.Error())

	// Ошибка: дата завершения в прошлом
	invalidTask = models.Task{
		Name: "Test Task",
		Due:  time.Now().Add(-24 * time.Hour),
	}
	_, err = service.Create(ctx, invalidTask)
	require.Error(t, err)
	require.Equal(t, "due date cannot be in the past", err.Error())
}

// Тест для GetByID
func TestTaskService_GetByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	ctx := context.Background()
	task := models.Task{
		ID:     1,
		Name:   "Test Task",
		Status: "Pending",
		Time:   time.Now(),
		Due:    time.Now().Add(24 * time.Hour),
	}

	// Успешный вызов
	mockRepo.On("GetById", ctx, task.ID).Return(&task, nil)

	result, err := service.GetByID(ctx, task.ID)

	require.NoError(t, err)
	require.Equal(t, &task, result)
	mockRepo.AssertExpectations(t)

	// Ошибка: задача не найдена
	mockRepo.On("GetById", ctx, 999).Return(nil, errors.New("task not found"))

	_, err = service.GetByID(ctx, 999)
	require.Error(t, err)
	require.Equal(t, "task not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// Тест для GetAll
func TestTaskService_GetAll(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	ctx := context.Background()
	tasks := []models.Task{
		{ID: 1, Name: "Task 1", Status: "Pending"},
		{ID: 2, Name: "Task 2", Status: "Completed"},
	}

	// Успешный вызов
	mockRepo.On("GetAll", ctx).Return(tasks, nil)

	result, err := service.GetAll(ctx)

	require.NoError(t, err)
	require.Equal(t, tasks, result)
	mockRepo.AssertExpectations(t)

	// Обработка ошибки
	mockRepo.On("GetAll", ctx).Return(nil, errors.New("no tasks found"))

	_, err = service.GetAll(ctx)
	require.Error(t, err)
	require.Equal(t, "no tasks found", err.Error())
	mockRepo.AssertExpectations(t)
}

// Тест для Update
func TestTaskService_Update(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	ctx := context.Background()
	existingTask := models.Task{
		ID:     1,
		Name:   "Task 1",
		Status: "Pending",
		Time:   time.Now(),
		Due:    time.Now().Add(24 * time.Hour),
	}
	updatedTask := models.Task{
		ID:     1,
		Name:   "Updated Task",
		Status: "Completed",
		Time:   existingTask.Time,
		Due:    existingTask.Due,
	}

	// Успешное обновление
	mockRepo.On("GetById", ctx, existingTask.ID).Return(&existingTask, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(task *models.Task) bool {
		return task.ID == updatedTask.ID && task.Name == updatedTask.Name &&
			task.Status == updatedTask.Status && task.Time.Equal(updatedTask.Time) &&
			task.Due.Equal(updatedTask.Due)
	})).Return(&updatedTask, nil)

	result, err := service.Update(ctx, updatedTask)

	require.NoError(t, err)
	require.Equal(t, updatedTask, result)
	mockRepo.AssertExpectations(t)

	// Ошибка: задача не найдена
	mockRepo.On("GetById", ctx, 999).Return(nil, errors.New("task not found"))

	_, err = service.Update(ctx, models.Task{ID: 999})
	require.Error(t, err)
	require.Equal(t, "task not found", err.Error())
	mockRepo.AssertExpectations(t)

	// Ошибка: имя задачи отсутствует
	_, err = service.Update(ctx, models.Task{ID: 1, Name: ""})
	require.Error(t, err)
	require.Equal(t, "task name is required", err.Error())

	// Ошибка: имя задачи слишком длинное
	longName := "This is a very long task name that exceeds fifty characters"
	_, err = service.Update(ctx, models.Task{ID: 1, Name: longName})
	require.Error(t, err)
	require.Equal(t, "task name must not exceed 50 characters", err.Error())

	// Ошибка: дата завершения в прошлом
	_, err = service.Update(ctx, models.Task{
		ID:   1,
		Name: "Valid Task",
		Due:  time.Now().Add(-1 * time.Hour),
	})
	require.Error(t, err)
	require.Equal(t, "due date cannot be in the past", err.Error())

	// Успешное обновление: установка значений по умолчанию
	partialUpdateTask := models.Task{
		ID:   1,
		Name: "Partial Update",
	}
	mockRepo.On("GetById", ctx, existingTask.ID).Return(&existingTask, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(task *models.Task) bool {
		return task.ID == partialUpdateTask.ID && task.Name == partialUpdateTask.Name &&
			task.Status == existingTask.Status && task.Time.Equal(existingTask.Time) &&
			task.Due.Equal(existingTask.Due)
	})).Return(&partialUpdateTask, nil)

	result, err = service.Update(ctx, partialUpdateTask)

	require.NoError(t, err)
	require.Equal(t, partialUpdateTask.ID, result.ID)
	require.Equal(t, partialUpdateTask.Name, result.Name)
	require.Equal(t, existingTask.Status, result.Status)
	require.True(t, result.Time.Equal(existingTask.Time))
	require.True(t, result.Due.Equal(existingTask.Due))
	mockRepo.AssertExpectations(t)
}

// Тест для Delete
func TestTaskService_Delete(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	ctx := context.Background()
	taskID := 1

	// Успешное удаление
	mockRepo.On("Delete", ctx, taskID).Return(nil)

	err := service.Delete(ctx, taskID)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Ошибка: задача не найдена
	mockRepo.On("Delete", ctx, 999).Return(errors.New("task not found"))

	err = service.Delete(ctx, 999)
	require.Error(t, err)
	require.Equal(t, "task not found", err.Error())
	mockRepo.AssertExpectations(t)
}
