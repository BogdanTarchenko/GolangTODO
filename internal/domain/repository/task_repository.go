package repository

import (
	"errors"
	"todo/internal/domain/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(task *model.Task) error
	Update(task *model.Task) error
	Delete(id string) error
	FindByID(id string) (*model.Task, error)
	FindAll() ([]*model.Task, error)
}
