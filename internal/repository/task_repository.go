package repository

import (
	"errors"
	"task-manager-api/internal/model"
)

var ErrNotFound = errors.New("record not found")

type TaskRepository interface {
	Create(task model.Task) (model.Task, error)
	GetByID(id int) (model.Task, error)
	List() ([]model.Task, error)
	Update(id int, input model.UpdateTaskInput) (model.Task, error)
	Delete(id int) error
}

type InMemoryTaskRepository struct {
	tasks  []model.Task
	nextID int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks:  []model.Task{},
		nextID: 1,
	}
}

func (r *InMemoryTaskRepository) Create(task model.Task) (model.Task, error) {
	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *InMemoryTaskRepository) GetByID(id int) (model.Task, error) {
	for _, t := range r.tasks {
		if t.ID == id {
			return t, nil
		}
	}

	return model.Task{}, ErrNotFound
}

func (r *InMemoryTaskRepository) List() ([]model.Task, error) {
	return r.tasks, nil
}

func (r *InMemoryTaskRepository) Delete(id int) error {
	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}

func (r *InMemoryTaskRepository) Update(id int, input model.UpdateTaskInput) (model.Task, error) {
	for i, t := range r.tasks {
		if t.ID == id {
			if input.Title != nil {
				r.tasks[i].Title = *input.Title
			}
			if input.Done != nil {
				r.tasks[i].Done = *input.Done
			}

			return r.tasks[i], nil
		}
	}

	return model.Task{}, ErrNotFound
}
