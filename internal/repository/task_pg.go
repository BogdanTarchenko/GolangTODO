package repository

import (
	"database/sql"
	"todo/internal/domain/model"
	"todo/internal/domain/repository"
)

type TaskPgRepository struct {
	db *sql.DB
}

func NewTaskPgRepository(db *sql.DB) *TaskPgRepository {
	return &TaskPgRepository{db: db}
}

func (r *TaskPgRepository) Create(task *model.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, deadline, status, priority, created_at, updated_at, is_completed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Deadline,
		task.Status,
		task.Priority,
		task.CreatedAt,
		task.UpdatedAt,
		task.IsCompleted,
	)
	return err
}

func (r *TaskPgRepository) Update(task *model.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, deadline = $3, status = $4, priority = $5, updated_at = $6, is_completed = $7
		WHERE id = $8
	`
	res, err := r.db.Exec(
		query,
		task.Title,
		task.Description,
		task.Deadline,
		task.Status,
		task.Priority,
		task.UpdatedAt,
		task.IsCompleted,
		task.ID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrTaskNotFound
	}
	return nil
}

func (r *TaskPgRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrTaskNotFound
	}
	return nil
}

func (r *TaskPgRepository) FindByID(id string) (*model.Task, error) {
	query := `
		SELECT id, title, description, deadline, status, priority, created_at, updated_at, is_completed
		FROM tasks WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var task model.Task
	var description sql.NullString
	var deadline sql.NullTime
	var updatedAt sql.NullTime

	err := row.Scan(
		&task.ID,
		&task.Title,
		&description,
		&deadline,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&updatedAt,
		&task.IsCompleted,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}
	if description.Valid {
		task.Description = &description.String
	}
	if deadline.Valid {
		task.Deadline = &deadline.Time
	}
	if updatedAt.Valid {
		task.UpdatedAt = &updatedAt.Time
	}

	return &task, nil
}
