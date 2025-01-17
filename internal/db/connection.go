package db

import (
	"WebTasks/cmd/WebTasks/config"
	"database/sql"
	"fmt"
	"log"
)

func DB(config *config.Config) *sql.DB {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DBName, config.DB.SSLMode, config.DB.SearchPath,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("База данных недоступна:", err)
	}
	log.Println("Подключение к базе данных успешно установлено!")
	return db
}

/* Герация идшника
   func generateID() string {

     id := uuid.New()
     return id.String()
   }
*/
