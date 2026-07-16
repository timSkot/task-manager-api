package repository

import (
	"database/sql"

	"task-manager-api/internal/model"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(task model.Task) (model.Task, error) {
	query := `INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRow(query, task.Title, task.Done).Scan(&task.ID)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *PostgresTaskRepository) GetByID(id int) (model.Task, error) {
	query := `SELECT id, title, done FROM tasks WHERE id = $1`

	var task model.Task
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Task{}, ErrNotFound
		}
		return model.Task{}, err
	}

	return task, nil
}

func (r *PostgresTaskRepository) List() ([]model.Task, error) {
	query := `SELECT id, title, done FROM tasks`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) Update(id int, input model.UpdateTaskInput) (model.Task, error) {
	task, err := r.GetByID(id)
	if err != nil {
		return model.Task{}, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Done != nil {
		task.Done = *input.Done
	}

	query := `UPDATE tasks SET title = $1, done = $2 WHERE id = $3`
	_, err = r.db.Exec(query, task.Title, task.Done, id)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *PostgresTaskRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
