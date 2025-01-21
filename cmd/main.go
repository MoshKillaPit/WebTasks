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
	"strconv"
)

func main() {
	// Чтение конфигурации
	cfg, err := config.ViperConfig()
	if err != nil {
		log.Fatalf("Ошибка чтения конфигурации базы данных: %v", err)
	}

	// Подключение к базе данных
	database, err := db.DB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	// Создание репозиториев
	userRepo := repositories.NewUserRepo(database)
	taskRepo := repositories.RepositoryForTasks(database)

	// Создание сервисов
	userService := services.NewUserService(userRepo)
	taskService := services.NewTaskService(taskRepo)

	// Создание маршрутов
	router := mux.NewRouter()

	// Применение глобальных middleware
	router.Use(handlers.LoggerMiddleware)      // Логирование запросов
	router.Use(handlers.CORSHeadersMiddleware) // Добавление CORS-заголовков

	// Регистрация маршрутов
	handlers.RegisterUserRoutes(router, userService)
	handlers.RegisterTaskRoutes(router, taskService)

	// Запуск сервера
	serverAddress := cfg.Server.IP + ":" + strconv.Itoa(cfg.Server.Port)
	log.Printf("Сервер запущен на %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

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
