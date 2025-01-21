# Используем Go версии 1.22 для сборки приложения
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main.go

# Используем минимальный образ с актуальной glibc
FROM debian:bookworm-slim AS final

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл приложения
COPY --from=builder /app/main .

# Копируем конфигурационные файлы
COPY config /app/config

# Указываем порт приложения
EXPOSE 8080

# Указываем команду для запуска
CMD ["./main"]
