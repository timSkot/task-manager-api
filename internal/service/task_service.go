package service

import (
	"errors"
	"strings"

	"task-manager-api/internal/model"
	"task-manager-api/internal/repository"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrInvalidTitle = errors.New("title is required")

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(input model.Task) (model.Task, error) {
	if strings.TrimSpace(input.Title) == "" {
		return model.Task{}, ErrInvalidTitle
	}

	return s.repo.Create(input), nil
}

func (s *TaskService) GetByID(id int) (model.Task, error) {
	task, ok := s.repo.GetByID(id)
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (s *TaskService) List() []model.Task {
	return s.repo.List()
}

func (s *TaskService) Update(id int, input model.UpdateTaskInput) (model.Task, error) {
	task, ok := s.repo.Update(id, input)
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (s *TaskService) Delete(id int) error {
	ok := s.repo.Delete(id)
	if !ok {
		return ErrTaskNotFound
	}

	return nil
}
