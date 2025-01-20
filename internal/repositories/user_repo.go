package repositories

import (
	"WebTasks/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetAll() ([]models.User, error)
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user models.User) (models.User, error) {
	query := `
		INSERT INTO users (name, key)
		VALUES (:name, :key)
		RETURNING id, name, key;
	`
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdUser models.User
		if err := rows.StructScan(&createdUser); err != nil {
			return models.User{}, err
		}
		return createdUser, nil
	}
	return models.User{}, nil
}

func (r *UserRepo) GetAll() ([]models.User, error) {
	query := "SELECT id, name, key FROM users;"
	var users []models.User
	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}
	return users, nil
}
