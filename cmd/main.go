package main

import (
	"WebTasks/config"
	"WebTasks/internal/db"
	"WebTasks/internal/handlers"
	"WebTasks/internal/repositories"
	"WebTasks/internal/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	// Инициализация данных о бд
	config, err := config.ViperConfig()
	if err != nil {
		log.Fatal("Ошибка чтения файла данных бд", err)
	}
	database, err := db.DB(config)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	taskRepo := repositories.RepositoryForTasks(database)
	taskService := services.NewTaskService(taskRepo)

	router := mux.NewRouter()
	handlers.RegisterTaskRoutes(router, taskService)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	/*
		прим миграций
					if err := db.ApplyMigrations(database); err != nil {
				        log.Fatalf("Ошибка применения миграций: %v", err)
				    }
			откат миграций
		if err := db.RollbackMigrations(database); err != nil {
		    log.Fatalf("Ошибка отката миграций: %v", err)
		}
	*/
}
