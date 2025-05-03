package repository

import "todo/internal/domain/model"

type TaskRepository interface {
	Create(task *model.Task) error
	Update(task *model.Task) error
	Delete(id string) error
	FindById(id string) (*model.Task, error)
	FindAll() ([]*model.Task, error)
}
