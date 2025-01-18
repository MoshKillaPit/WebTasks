package db

import (
	"WebTasks/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

func DB(config *config.Config) (*sqlx.DB, error) {
	log.Printf("Подключение к базе данных с параметрами: host=%s port=%d user=%s dbname=%s sslmode=%s",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.DBName, config.DB.SSLMode)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DBName, config.DB.SSLMode, config.DB.SearchPath,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}

	log.Println("Подключение к базе данных успешно установлено!")
	return db, nil
}

/* Герация идшника
   func generateID() string {

     id := uuid.New()
     return id.String()
   }
*/
