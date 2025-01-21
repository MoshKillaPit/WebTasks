package main

import (
	"WebTasks/config"
	"WebTasks/internal/db"
	"WebTasks/internal/handlers"
	"WebTasks/internal/repositories"
	"WebTasks/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Чтение конфигурации
	cfg, err := config.ViperConfig()
	if err != nil {
		log.Printf("Ошибка чтения конфигурации базы данных: %v", err)
		return
	}

	// Подключение к базе данных
	database, err := db.DB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return
	}

	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	}()

	// Создание репозиториев
	userRepo := repositories.NewUserRepo(database)
	taskRepo := repositories.RepositoryForTasks(database)

	// Создание сервисов
	userService := services.NewUserService(userRepo)
	taskService := services.NewTaskService(taskRepo)

	// Создание обработчиков
	taskHandler := handlers.NewHandler(taskService)
	userHandler := handlers.NewUserHandler(userService)

	// Создание маршрутов
	router := mux.NewRouter()

	// Применение глобальных middleware
	router.Use(handlers.LoggerMiddleware) // Логирование запросов
	router.Use(handlers.AuthMiddleware)

	// Регистрация маршрутов
	handlers.RegisterUserRoutes(router, userHandler)
	handlers.RegisterTaskRoutes(router, taskHandler)

	// Запуск сервера
	serverAddress := cfg.Server.IP + ":" + strconv.Itoa(cfg.Server.Port)
	log.Printf("Сервер запущен на %s", serverAddress)

	err = http.ListenAndServe(serverAddress, router)
	if err != nil {
		log.Printf("Ошибка запуска сервера: %v", err)
		return
	}
}
