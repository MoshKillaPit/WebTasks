package services

import (
	"WebTasks/internal/models"
	"WebTasks/internal/repositories"
	"context"
	"errors"
	"time"
)

type TaskService interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
	GetByID(ctx context.Context, id int) (*models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, task models.Task) (models.Task, error)
	Delete(ctx context.Context, id int) error
}

type taskServiceImpl struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &taskServiceImpl{repo: repo}
}

func (s *taskServiceImpl) Create(ctx context.Context, task models.Task) (models.Task, error) {
	if task.Name == "" {
		return models.Task{}, errors.New("task name is required")
	}

	if len(task.Name) > 50 {
		return models.Task{}, errors.New("task name is too long")
	}

	if !task.Due.IsZero() && task.Due.Before(time.Now()) {
		return models.Task{}, errors.New("due date cannot be in the past")
	}

	createdTask, err := s.repo.Create(ctx, &task)
	if err != nil {
		return models.Task{}, err
	}

	return *createdTask, nil
}

func (s *taskServiceImpl) GetAll(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *taskServiceImpl) GetByID(ctx context.Context, id int) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskServiceImpl) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *taskServiceImpl) Update(ctx context.Context, task models.Task) (models.Task, error) {
	existingTask, err := s.repo.GetByID(ctx, task.ID)
	if err != nil {
		return models.Task{}, err
	}

	if task.Name == "" {
		return models.Task{}, errors.New("task name is required")
	}

	if len(task.Name) > 50 {
		return models.Task{}, errors.New("task name must not exceed 50 characters")
	}

	if !task.Due.IsZero() && task.Due.Before(time.Now()) {
		return models.Task{}, errors.New("due date cannot be in the past")
	}

	if task.Status == "" {
		task.Status = existingTask.Status
	}

	if task.Time.IsZero() {
		task.Time = existingTask.Time
	}

	if task.Due.IsZero() {
		task.Due = existingTask.Due
	}

	updatedTask, err := s.repo.Update(ctx, &task)
	if err != nil {
		return models.Task{}, err
	}

	return *updatedTask, nil
}
