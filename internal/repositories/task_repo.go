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
	defer rows.Close()

	if rows.Next() {
		var createdTask models.Task
		if err := rows.StructScan(&createdTask); err != nil {
			log.Printf("Error scanning created task: %v", err)
			return nil, err
		}
		return &createdTask, nil
	}

	log.Printf("Task creation failed, no rows returned")
	return nil, nil
}

func (r *TaskRepo) GetById(ctx context.Context, id int) (*models.Task, error) {
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
	defer rows.Close()

	if rows.Next() {
		var updatedTask models.Task
		if err := rows.StructScan(&updatedTask); err != nil {
			log.Printf("Error scanning updated task: %v", err)
			return nil, err
		}
		return &updatedTask, nil
	}

	log.Printf("Task update failed, no rows returned")
	return nil, nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, DeleteTaskQuery, id)
	if err != nil {
		log.Printf("Error executing DeleteTaskQuery for id %d: %v", id, err)
		return err
	}
	return nil
}
