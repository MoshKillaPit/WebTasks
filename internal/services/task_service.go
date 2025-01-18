package services

import (
	"WebTasks/internal/models"
	"WebTasks/internal/repositories"
	"database/sql"
	"errors"
	"time"
)

type TaskService interface {
	Create(task models.Task) (models.Task, error)
	GetByID(id int) (*models.Task, error) // Меняем models.Task на *models.Task
	GetAll() ([]models.Task, error)
	Update(task models.Task) (models.Task, error)
	Delete(id int) error
}

type TaskServiceImpl struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{repo: repo}
}

func (s *TaskServiceImpl) Create(task models.Task) (models.Task, error) {
	if task.Name == "" {
		return models.Task{}, errors.New("task name is required")
	}
	if len(task.Name) > 50 {
		return models.Task{}, errors.New("task name is too long")
	}
	if !task.Due.IsZero() && task.Due.Before(time.Now()) {
		return models.Task{}, errors.New("due date cannot be in the past")
	}

	createdTask, err := s.repo.Create(&task)
	if err != nil {
		return models.Task{}, err
	}

	return *createdTask, nil // Разыменовываем указатель для возврата значения
}

func (s *TaskServiceImpl) GetAll() ([]models.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskServiceImpl) GetByID(id int) (*models.Task, error) {
	return s.repo.GetById(id)
}

func (s *TaskServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TaskServiceImpl) Update(task models.Task) (models.Task, error) {
	existingTask, err := s.repo.GetById(task.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, errors.New("task not found")
		}
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

	updatedTask, err := s.repo.Update(&task)
	if err != nil {
		return models.Task{}, err
	}

	// Разыменование указателя для возврата
	return *updatedTask, nil
}
