package services

import (
	"WebTasks/internal/models"
	"context"
	"errors"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetByID(ctx context.Context, id int) (*models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) (*models.Task, error)
	Delete(ctx context.Context, id int) error
}

type taskService struct {
	repo TaskRepository
}

// Create добавляет новую задачу в хранилище
func (s *taskService) Create(ctx context.Context, task models.Task) (models.Task, error) {
	// Проверяем обязательные поля
	if task.Name == "" {
		return models.Task{}, errors.New("task name is required")
	}
	// В тестах для Create ожидается ошибка "task name is too long"
	// если имя длиннее 50 символов
	if len(task.Name) > 50 {
		return models.Task{}, errors.New("task name is too long")
	}
	if !task.Due.IsZero() && task.Due.Before(time.Now()) {
		return models.Task{}, errors.New("due date cannot be in the past")
	}

	created, err := s.repo.Create(ctx, &task)
	if err != nil {
		return models.Task{}, err
	}
	return *created, nil
}

// GetByID возвращает задачу по её ID
func (s *taskService) GetByID(ctx context.Context, id int) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAll возвращает список всех задач
func (s *taskService) GetAll(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		// Ранее здесь могли проглатывать ошибку;
		// теперь возвращаем её в вызывающий код:
		return nil, err
	}
	return tasks, nil
}

// Update обновляет существующую задачу
func (s *taskService) Update(ctx context.Context, newData models.Task) (models.Task, error) {
	// Сначала нужно получить существующую задачу, чтобы понять, что у неё уже есть
	existing, err := s.repo.GetByID(ctx, newData.ID)
	if err != nil {
		return models.Task{}, err
	}

	if newData.Name == "" {
		return models.Task{}, errors.New("task name is required")
	}
	// В тестах для Update ожидается ошибка "task name must not exceed 50 characters"
	if len(newData.Name) > 50 {
		return models.Task{}, errors.New("task name must not exceed 50 characters")
	}
	if !newData.Due.IsZero() && newData.Due.Before(time.Now()) {
		return models.Task{}, errors.New("due date cannot be in the past")
	}

	// Частичное обновление: перенимаем некоторые поля из существующей задачи,
	// если в newData они не заданы или равны нулевым значениям.
	finalTask := *existing

	// Имя меняем всегда, т.к. тесты требуют ошибка при пустом имени
	finalTask.Name = newData.Name

	// Если Status не передан (пустая строка), оставляем старый
	if newData.Status != "" {
		finalTask.Status = newData.Status
	}
	// Если Time не передан (zero time), оставляем старый
	if !newData.Time.IsZero() {
		finalTask.Time = newData.Time
	}
	// Если Due не передан (zero time), оставляем старый
	if !newData.Due.IsZero() {
		finalTask.Due = newData.Due
	}

	updated, err := s.repo.Update(ctx, &finalTask)
	if err != nil {
		return models.Task{}, err
	}
	return *updated, nil
}

// Delete удаляет задачу по ID
func (s *taskService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
