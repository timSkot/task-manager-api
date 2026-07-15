package repository

import "task-manager-api/internal/model"

type TaskRepository interface {
	Create(task model.Task) model.Task
	GetByID(id int) (model.Task, bool)
	List() []model.Task
	Update(id int, input model.UpdateTaskInput) (model.Task, bool)
	Delete(id int) bool
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

func (r *InMemoryTaskRepository) Create(task model.Task) model.Task {
	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, task)

	return task
}

func (r *InMemoryTaskRepository) GetByID(id int) (model.Task, bool) {
	for _, t := range r.tasks {
		if t.ID == id {
			return t, true
		}
	}

	return model.Task{}, false
}

func (r *InMemoryTaskRepository) List() []model.Task {
	return r.tasks
}

func (r *InMemoryTaskRepository) Delete(id int) bool {
	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return true
		}
	}

	return false
}

func (r *InMemoryTaskRepository) Update(id int, input model.UpdateTaskInput) (model.Task, bool) {
	for i, t := range r.tasks {
		if t.ID == id {
			if input.Title != nil {
				r.tasks[i].Title = *input.Title
			}
			if input.Done != nil {
				r.tasks[i].Done = *input.Done
			}

			return r.tasks[i], true
		}
	}

	return model.Task{}, false
}
