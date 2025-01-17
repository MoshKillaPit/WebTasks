package main

import (
	config2 "WebTasks/cmd/WebTasks/config"
	"WebTasks/cmd/WebTasks/internal/db"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// Инициализация данных о бд
	config, err := config2.ViperConfig()
	if err != nil {
		log.Fatal("Ошибка чтения файла данных бд", err)
	}
	db.DB(config)

}
