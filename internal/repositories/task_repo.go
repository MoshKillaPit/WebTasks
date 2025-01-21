package repositories

import (
	"WebTasks/internal/models"
	"context"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetByID(ctx context.Context, id int) (*models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) (*models.Task, error)
	Delete(ctx context.Context, id int) error
}

type TaskRepo struct {
	db *sqlx.DB
}

func RepositoryForTasks(db *sqlx.DB) TaskRepository {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	rows, err := r.db.NamedQueryContext(ctx, CreateTaskQuery, task)
	if err != nil {
		log.Printf("Error executing CreateTaskQuery: %v", err)
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Printf("Error closing rows in Create: %v", closeErr)
		}
	}()

	var createdTask models.Task
	if rows.Next() {
		if err := rows.StructScan(&createdTask); err != nil {
			log.Printf("Error scanning created task: %v", err)
			return nil, err
		}

		return &createdTask, nil
	}

	log.Printf("Task creation failed, no rows returned")

	// Пустая строка перед return
	return nil, errors.New("task creation failed: no rows returned")
}

func (r *TaskRepo) GetByID(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task

	err := r.db.GetContext(ctx, &task, GetTaskByIDQuery, id)
	if err != nil {
		log.Printf("Error executing GetTaskByIDQuery for id %d: %v", id, err)
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task

	err := r.db.SelectContext(ctx, &tasks, GetAllTasksQuery)
	if err != nil {
		log.Printf("Error executing GetAllTasksQuery: %v", err)
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepo) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	rows, err := r.db.NamedQueryContext(ctx, UpdateTaskQuery, task)
	if err != nil {
		log.Printf("Error executing UpdateTaskQuery: %v", err)
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Printf("Error closing rows in Update: %v", closeErr)
		}
	}()

	var updatedTask models.Task
	if rows.Next() {
		if err := rows.StructScan(&updatedTask); err != nil {
			log.Printf("Error scanning updated task: %v", err)
			return nil, err
		}

		return &updatedTask, nil
	}

	log.Printf("Task update failed, no rows returned")

	// Пустая строка перед return
	return nil, errors.New("task update failed: no rows returned")
}

func (r *TaskRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, DeleteTaskQuery, id)
	if err != nil {
		log.Printf("Error executing DeleteTaskQuery for id %d: %v", id, err)
		return err
	}

	return nil
}
