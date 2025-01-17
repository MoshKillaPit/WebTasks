package db

import (
	"WebTasks/cmd/WebTasks/config"
	"database/sql"
	"fmt"
	"log"
)

func DB() *sql.DB {
	cfg, err := config.ViperConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode, cfg.DB.SearchPath,
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
