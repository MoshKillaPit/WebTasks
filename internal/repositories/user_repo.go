package repositories

import (
	"WebTasks/internal/models"
	"context"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int) error
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	rows, err := r.db.NamedQueryContext(ctx, CreateUserQuery, user)
	if err != nil {
		logError("CreateUserQuery", err)
		return nil, err
	}

	// Отложенный вызов с проверкой возможной ошибки закрытия
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logError("Close rows in Create", closeErr)
		}
	}()

	if rows.Next() {
		var createdUser models.User
		if err := rows.StructScan(&createdUser); err != nil {
			logError("StructScan (Create)", err)
			return nil, err
		}

		return &createdUser, nil
	}

	log.Printf("User creation failed: no rows returned")

	return nil, errors.New("user creation failed")
}

func (r *UserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.SelectContext(ctx, &users, GetAllUsersQuery)
	if err != nil {
		logError("GetAllUsersQuery", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, GetUserByIDQuery, id)
	if err != nil {
		logError("GetUserByIDQuery", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	rows, err := r.db.NamedQueryContext(ctx, UpdateUserQuery, user)
	if err != nil {
		logError("UpdateUserQuery", err)
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logError("Close rows in Update", closeErr)
		}
	}()

	if rows.Next() {
		var updatedUser models.User
		if err := rows.StructScan(&updatedUser); err != nil {
			logError("StructScan (Update)", err)
			return nil, err
		}

		return &updatedUser, nil
	}

	log.Printf("User update failed: no rows returned")

	return nil, errors.New("user update failed")
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, DeleteUserQuery, id)
	if err != nil {
		logError("DeleteUserQuery", err)
		return err
	}

	return nil
}

func logError(query string, err error) {
	log.Printf("Error executing query '%s': %v", query, err)
}
