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

	task, err := s.repo.Create(input)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (s *TaskService) GetByID(id int) (model.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Task{}, ErrTaskNotFound
		}
		return model.Task{}, err
	}

	return task, nil
}

func (s *TaskService) List() ([]model.Task, error) {
	tasks, err := s.repo.List()

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) Update(id int, input model.UpdateTaskInput) (model.Task, error) {
	task, err := s.repo.Update(id, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Task{}, ErrTaskNotFound
		}
		return model.Task{}, err
	}

	return task, nil
}

func (s *TaskService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	return nil
}
