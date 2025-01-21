package handlers_test

import (
	"WebTasks/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	// Создаем тестовый обработчик, который возвращает HTTP 200
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond) // Симулируем задержку для проверки длительности
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем тестовый обработчик в LoggerMiddleware
	loggerMiddleware := handlers.LoggerMiddleware(nextHandler)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	loggerMiddleware.ServeHTTP(rr, req)

	// Проверяем, что статус-код 200
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_ValidHeader(t *testing.T) {
	// Создаем тестовый обработчик, который возвращает HTTP 200
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем тестовый обработчик в AuthMiddleware
	authMiddleware := handlers.AuthMiddleware(nextHandler)

	// Создаем тестовый HTTP-запрос с заголовком Authorization
	req := httptest.NewRequest(http.MethodGet, "/auth-test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	authMiddleware.ServeHTTP(rr, req)

	// Проверяем, что статус-код 200
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	// Создаем тестовый обработчик, который возвращает HTTP 200
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем тестовый обработчик в AuthMiddleware
	authMiddleware := handlers.AuthMiddleware(nextHandler)

	// Создаем тестовый HTTP-запрос без заголовка Authorization
	req := httptest.NewRequest(http.MethodGet, "/auth-test", nil)
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	authMiddleware.ServeHTTP(rr, req)

	// Проверяем, что статус-код 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Unauthorized: Missing Authorization Header")
}
