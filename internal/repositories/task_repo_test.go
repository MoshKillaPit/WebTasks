package repositories_test

import (
	"WebTasks/internal/models"
	"WebTasks/internal/repositories"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.RepositoryForTasks(sqlxDB)

	task := &models.Task{
		Name:   "Test Task",
		Status: "Pending",
		Time:   time.Now(),
		Due:    time.Now().Add(24 * time.Hour),
		UserID: 1,
	}

	mock.ExpectQuery(`INSERT INTO public.tasks`).
		WithArgs(task.Name, task.Status, task.Time, task.Due, task.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "time", "due", "user_id"}).
			AddRow(1, "Test Task", "Pending", task.Time, task.Due, 1))

	ctx := context.Background()
	createdTask, err := repo.Create(ctx, task)

	assert.NoError(t, err)
	assert.NotNil(t, createdTask)
	assert.Equal(t, 1, createdTask.ID)
	assert.Equal(t, "Test Task", createdTask.Name)
	assert.Equal(t, "Pending", createdTask.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaskRepo_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.RepositoryForTasks(sqlxDB)

	mock.ExpectQuery(`SELECT id, name, status, time, due FROM public.tasks`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "time", "due"}).
			AddRow(1, "Task 1", "Pending", time.Now(), time.Now().Add(24*time.Hour)).
			AddRow(2, "Task 2", "Completed", time.Now(), time.Now().Add(48*time.Hour)))

	ctx := context.Background()
	tasks, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "Task 1", tasks[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaskRepo_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.RepositoryForTasks(sqlxDB)

	mock.ExpectQuery(`SELECT id, name, status, time, due FROM public.tasks WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "time", "due"}).
			AddRow(1, "Task 1", "Pending", time.Now(), time.Now().Add(24*time.Hour)))

	ctx := context.Background()
	task, err := repo.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "Task 1", task.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaskRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.RepositoryForTasks(sqlxDB)

	task := &models.Task{
		ID:     1,
		Name:   "Updated Task",
		Status: "Completed",
		Time:   time.Now(),
		Due:    time.Now().Add(48 * time.Hour),
	}

	mock.ExpectQuery(`UPDATE public.tasks SET`).
		WithArgs(task.Name, task.Status, task.Time, task.Due, task.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "time", "due"}).
			AddRow(1, "Updated Task", "Completed", task.Time, task.Due))

	ctx := context.Background()
	updatedTask, err := repo.Update(ctx, task)

	assert.NoError(t, err)
	assert.NotNil(t, updatedTask)
	assert.Equal(t, 1, updatedTask.ID)
	assert.Equal(t, "Updated Task", updatedTask.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaskRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.RepositoryForTasks(sqlxDB)

	mock.ExpectExec(`DELETE FROM public.tasks WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err = repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
