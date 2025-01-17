package main

import (
	"WebTasks/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к бд
	db.DB()
}
