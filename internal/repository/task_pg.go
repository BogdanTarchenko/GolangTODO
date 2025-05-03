package repository

import (
	"database/sql"
	"todo/internal/domain/model"
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
