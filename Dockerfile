# Dockerfile
# Билд-стадия
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev
# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod ./
COPY go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/contact-api ./cmd/contact-api

# Релиз-стадия
FROM alpine AS runner

WORKDIR /app

# Копируем собранное приложение из билд-стадии
COPY --from=builder /app/contact-api /app/contact-api

# Копируем конфигурационные файлы
COPY --from=builder /app/config /app/config

COPY --from=builder /app/.env ./.env

# Открываем порт
EXPOSE 8080

# Запуск приложения
CMD ["/app/contact-api"]
