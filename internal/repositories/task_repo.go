package repositories

import (
	"WebTasks/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type TaskRepository interface {
	Create(task *models.Task) (*models.Task, error)
	GetById(id int) (*models.Task, error)
	GetAll() ([]models.Task, error)
	Update(task *models.Task) (*models.Task, error)
	Delete(id int) error
}

type Task struct {
	db *sqlx.DB
}

func RepositoryForTasks(db *sqlx.DB) *Task {
	return &Task{db: db}
}

func (r *Task) Create(task *models.Task) (*models.Task, error) {
	rows, err := r.db.NamedQuery(CreateTaskQuery, task)
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

func (r *Task) GetById(id int) (*models.Task, error) {
	var task models.Task
	err := r.db.Get(&task, GetTaskByIDQuery, id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Task) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Select(&tasks, GetAllTasksQuery)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Task) Update(task *models.Task) (*models.Task, error) {
	rows, err := r.db.NamedQuery(UpdateTaskQuery, task)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var updatedTask models.Task
		err := rows.StructScan(&updatedTask)
		if err != nil {
			return nil, err
		}
		return &updatedTask, nil
	}
	return nil, nil
}

func (r *Task) Delete(id int) error {
	_, err := r.db.Exec(DeleteTaskQuery, id)
	return err
}
