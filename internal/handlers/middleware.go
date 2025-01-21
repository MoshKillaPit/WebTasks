package handlers

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware логирует каждый HTTP-запрос: метод, URL, длительность выполнения
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Начало обработки запроса: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r) // Передача управления следующему обработчику

		duration := time.Since(start)
		log.Printf("Запрос обработан: %s %s за %v", r.Method, r.URL.Path, duration)
	})
}

// AuthMiddleware проверяет наличие заголовка Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r) // Передача управления следующему обработчику
	})
}
