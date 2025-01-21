package repositories

import (
	"WebTasks/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetById(ctx context.Context, id int) (*models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) (*models.Task, error)
	Delete(ctx context.Context, id int) error
}

type Task struct {
	db *sqlx.DB
}

func RepositoryForTasks(db *sqlx.DB) *Task {
	return &Task{db: db}
}

func (r *Task) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	rows, err := r.db.NamedQueryContext(ctx, CreateTaskQuery, task)
	if err != nil {
		log.Printf("Error executing CreateTaskQuery: %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdTask models.Task
		err := rows.StructScan(&createdTask)
		if err != nil {
			log.Printf("Error scanning created task: %v", err)
			return nil, err
		}
		return &createdTask, nil
	}
	return nil, nil
}

func (r *Task) GetById(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task
	err := r.db.GetContext(ctx, &task, GetTaskByIDQuery, id)
	if err != nil {
		log.Printf("Error executing GetTaskByIDQuery: %v", err)
		return nil, err
	}
	return &task, nil
}

func (r *Task) GetAll(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.SelectContext(ctx, &tasks, GetAllTasksQuery)
	if err != nil {
		log.Printf("Error executing GetAllTasksQuery: %v", err)
		return nil, err
	}
	return tasks, nil
}

func (r *Task) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	rows, err := r.db.NamedQueryContext(ctx, UpdateTaskQuery, task)
	if err != nil {
		log.Printf("Error executing UpdateTaskQuery: %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var updatedTask models.Task
		err := rows.StructScan(&updatedTask)
		if err != nil {
			log.Printf("Error scanning updated task: %v", err)
			return nil, err
		}
		return &updatedTask, nil
	}
	return nil, nil
}

func (r *Task) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, DeleteTaskQuery, id)
	if err != nil {
		log.Printf("Error executing DeleteTaskQuery: %v", err)
	}
	return err
}
