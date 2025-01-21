package repositories_test

import (
	"WebTasks/internal/models"
	"WebTasks/internal/repositories"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepo(sqlxDB)

	user := &models.User{
		Name: "Test User",
		Key:  "test-key",
	}

	mock.ExpectQuery(`INSERT INTO public.users`).
		WithArgs(user.Name, user.Key).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "key"}).AddRow(1, "Test User", "test-key"))

	ctx := context.Background()
	createdUser, err := repo.Create(ctx, user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, 1, createdUser.ID)
	assert.Equal(t, "Test User", createdUser.Name)
	assert.Equal(t, "test-key", createdUser.Key)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepo(sqlxDB)

	mock.ExpectQuery(`SELECT id, name, key FROM public.users`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "key"}).
			AddRow(1, "User1", "key1").
			AddRow(2, "User2", "key2"))

	ctx := context.Background()
	users, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "User1", users[0].Name)
	assert.Equal(t, "key1", users[0].Key)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepo(sqlxDB)

	mock.ExpectQuery(`SELECT id, name, key FROM public.users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "key"}).AddRow(1, "User1", "key1"))

	ctx := context.Background()
	user, err := repo.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "User1", user.Name)
	assert.Equal(t, "key1", user.Key)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepo(sqlxDB)

	user := &models.User{
		ID:   1,
		Name: "Updated User",
		Key:  "updated-key",
	}

	mock.ExpectQuery(`UPDATE public.users SET`).
		WithArgs(user.Name, user.Key, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "key"}).AddRow(1, "Updated User", "updated-key"))

	ctx := context.Background()
	updatedUser, err := repo.Update(ctx, user)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, 1, updatedUser.ID)
	assert.Equal(t, "Updated User", updatedUser.Name)
	assert.Equal(t, "updated-key", updatedUser.Key)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repositories.NewUserRepo(sqlxDB)

	mock.ExpectExec(`DELETE FROM public.users WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err = repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
