package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"task-manager-api/internal/model"
	"task-manager-api/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input model.Task
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid JSON:", err)
		return
	}

	task, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTitle) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	// вызвать h.service.List(), там нет error вообще — просто вернуть JSON
	tasks := h.service.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// r.PathValue("id") → strconv.Atoi → h.service.GetByID(id)
	// если err — errors.Is(err, service.ErrTaskNotFound) → 404
	// иначе — вернуть JSON
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid id:", err)
		return
	}

	task, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid id:", err)
		return
	}

	var input model.UpdateTaskInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid JSON:", err)
		return
	}

	task, updateErr := h.service.Update(id, input) // не игнорируем первое значение
	if updateErr != nil {
		if errors.Is(updateErr, service.ErrTaskNotFound) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, updateErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, updateErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// id из пути → h.service.Delete(id)
	// если err — 404
	// если nil — 204 No Content, без тела
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid id:", err)
		return
	}

	deleteErr := h.service.Delete(id)
	if deleteErr != nil {
		if errors.Is(deleteErr, service.ErrTaskNotFound) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, deleteErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
