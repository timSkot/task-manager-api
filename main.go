package main

import (
	"fmt"
	"net/http"

	"task-manager-api/internal/handler"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/service"
)

func main() {
	repo := repository.NewInMemoryTaskRepository()
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	http.HandleFunc("GET /tasks", h.ListTasks)
	http.HandleFunc("POST /tasks", h.CreateTask)
	http.HandleFunc("GET /tasks/{id}", h.GetTask)
	http.HandleFunc("PATCH /tasks/{id}", h.UpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)

	fmt.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
