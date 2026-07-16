# Стадия 1: сборка бинарника
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o task-manager-api main.go

# Стадия 2: минимальный финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/task-manager-api .

EXPOSE 8080

CMD ["./task-manager-api"]