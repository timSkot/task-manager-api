package service

import (
	"errors"
	"task-manager-api/internal/model"
	"task-manager-api/internal/repository"
	"testing"
)

type MockTaskRepository struct {
	CreateFunc  func(task model.Task) (model.Task, error)
	GetByIDFunc func(id int) (model.Task, error)
	ListFunc    func() ([]model.Task, error)
	UpdateFunc  func(id int, input model.UpdateTaskInput) (model.Task, error)
	DeleteFunc  func(id int) error
}

func (m *MockTaskRepository) Create(task model.Task) (model.Task, error) {
	return m.CreateFunc(task)
}

func (m *MockTaskRepository) GetByID(id int) (model.Task, error) {
	return m.GetByIDFunc(id)
}

func (m *MockTaskRepository) List() ([]model.Task, error) {
	return m.ListFunc()
}

func (m *MockTaskRepository) Update(id int, input model.UpdateTaskInput) (model.Task, error) {
	return m.UpdateFunc(id, input)
}

func (m *MockTaskRepository) Delete(id int) error {
	return m.DeleteFunc(id)
}

func TestTaskService_Create_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		CreateFunc: func(task model.Task) (model.Task, error) {
			task.ID = 1
			return task, nil
		},
	}

	svc := NewTaskService(mockRepo)

	result, err := svc.Create(model.Task{Title: "Learn Go"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.Title != "Learn Go" {
		t.Errorf("expected title 'Learn Go', got %q", result.Title)
	}
}

func TestTaskService_Create_EmptyTitle(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	svc := NewTaskService(mockRepo)

	_, err := svc.Create(model.Task{Title: ""})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrInvalidTitle) {
		t.Errorf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestTaskService_Create_RepositoryError(t *testing.T) {
	someErr := errors.New("database connection failed") // симулируем "сбой базы"

	mockRepo := &MockTaskRepository{
		CreateFunc: func(task model.Task) (model.Task, error) {
			return model.Task{}, someErr
		},
	}

	svc := NewTaskService(mockRepo)

	_, err := svc.Create(model.Task{Title: "Learn Go"})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, someErr) {
		t.Errorf("expected someErr, got %v", err)
	}
}

func TestTaskService_GetByID_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetByIDFunc: func(id int) (model.Task, error) {
			return model.Task{ID: id, Title: "Learn Go"}, nil
		},
	}

	svc := NewTaskService(mockRepo)

	result, err := svc.GetByID(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestTaskService_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetByIDFunc: func(id int) (model.Task, error) {
			return model.Task{}, repository.ErrNotFound
		},
	}

	svc := NewTaskService(mockRepo)

	_, err := svc.GetByID(1)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskService_List_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		ListFunc: func() ([]model.Task, error) {
			return []model.Task{
				{ID: 1, Title: "Learn Go", Done: false},
				{ID: 2, Title: "Buy milk", Done: true},
			}, nil
		},
	}

	svc := NewTaskService(mockRepo)

	result, err := svc.List()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(result))
	}
}

func TestTaskService_List_Error(t *testing.T) {
	someErr := errors.New("database connection failed")

	mockRepo := &MockTaskRepository{
		ListFunc: func() ([]model.Task, error) {
			return []model.Task{}, someErr
		},
	}

	svc := NewTaskService(mockRepo)

	_, err := svc.List()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, someErr) {
		t.Errorf("expected someErr, got %v", err)
	}
}

func TestTaskService_Update_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		UpdateFunc: func(id int, input model.UpdateTaskInput) (model.Task, error) {
			return model.Task{ID: id, Title: "Updated title", Done: true}, nil
		},
	}

	svc := NewTaskService(mockRepo)

	newTitle := "Updated title"
	input := model.UpdateTaskInput{Title: &newTitle}

	result, err := svc.Update(1, input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Title != "Updated title" {
		t.Errorf("expected title 'Updated title', got %q", result.Title)
	}
}

func TestTaskService_Update_NotFound(t *testing.T) {
	mockRepo := &MockTaskRepository{
		UpdateFunc: func(id int, input model.UpdateTaskInput) (model.Task, error) {
			return model.Task{}, repository.ErrNotFound
		},
	}

	svc := NewTaskService(mockRepo)

	newTitle := "Updated title"
	input := model.UpdateTaskInput{Title: &newTitle}

	_, err := svc.Update(1, input)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskService_Delete_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteFunc: func(id int) error {
			return nil
		},
	}

	svc := NewTaskService(mockRepo)

	err := svc.Delete(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTaskService_Delete_NotFound(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteFunc: func(id int) error {
			return repository.ErrNotFound
		},
	}

	svc := NewTaskService(mockRepo)

	err := svc.Delete(1)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}
