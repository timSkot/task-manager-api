package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"task-manager-api/internal/handler"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database!")

	repo := repository.NewPostgresTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	http.HandleFunc("GET /tasks", h.ListTasks)
	http.HandleFunc("POST /tasks", h.CreateTask)
	http.HandleFunc("GET /tasks/{id}", h.GetTask)
	http.HandleFunc("PATCH /tasks/{id}", h.UpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
	http.Handle("GET /metrics", promhttp.Handler())
	wrappedMux := handler.LoggingMiddleware(http.DefaultServeMux)

	fmt.Println("Server starting on :8080...")
	err = http.ListenAndServe(":8080", wrappedMux)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
