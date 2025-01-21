package repositories

import (
	"WebTasks/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetAll() ([]models.User, error)
	GetById(id int) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int) error
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user models.User) (models.User, error) {
	rows, err := r.db.NamedQuery(CreateUserQuery, user)
	if err != nil {
		log.Printf("Error executing CreateUserQuery: %v", err)
		return models.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdUser models.User
		err := rows.StructScan(&createdUser)
		if err != nil {
			log.Printf("Error scanning created user: %v", err)
			return models.User{}, err
		}
		return createdUser, nil
	}
	return models.User{}, nil
}

func (r *UserRepo) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, GetAllUsersQuery)
	if err != nil {
		log.Printf("Error executing GetAllUsersQuery: %v", err)
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) GetById(id int) (models.User, error) {
	var user models.User
	err := r.db.Get(&user, GetUserByIDQuery, id)
	if err != nil {
		log.Printf("Error executing GetUserByIDQuery: %v", err)
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepo) Update(user models.User) (models.User, error) {
	rows, err := r.db.NamedQuery(UpdateUserQuery, user)
	if err != nil {
		log.Printf("Error executing UpdateUserQuery: %v", err)
		return models.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var updatedUser models.User
		err := rows.StructScan(&updatedUser)
		if err != nil {
			log.Printf("Error scanning updated user: %v", err)
			return models.User{}, err
		}
		return updatedUser, nil
	}
	return models.User{}, nil
}

func (r *UserRepo) Delete(id int) error {
	_, err := r.db.Exec(DeleteUserQuery, id)
	if err != nil {
		log.Printf("Error executing DeleteUserQuery: %v", err)
	}
	return err
}
