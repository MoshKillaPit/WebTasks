package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func ApplyMigrations(db *sqlx.DB) error {
	queries := []string{
		// Создание таблицы tasks
		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			due TIMESTAMP DEFAULT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			key VARCHAR(255) NOT NULL UNIQUE
		);`,

		`ALTER TABLE tasks ADD COLUMN IF NOT EXISTS user_id INT REFERENCES users(id) ON DELETE CASCADE;`,
	}

	// Выполнение миграций
	for _, query := range queries {
		log.Printf("Выполнение миграции: %s", query)

		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Ошибка выполнения миграции: %v", err)

			return err
		}
	}

	log.Println("Миграции успешно применены!")

	return nil
}

func RollbackMigrations(db *sqlx.DB) error {
	queries := []string{
		`DROP TABLE IF EXISTS tasks;`,
		`DROP TABLE IF EXISTS users;`,
	}

	for _, query := range queries {
		log.Printf("Откат миграции: %s", query)

		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Ошибка отката миграции: %v", err)

			return err
		}
	}

	log.Println("Откат миграций успешно выполнен!")

	return nil
}
