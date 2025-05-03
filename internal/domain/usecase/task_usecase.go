package usecase

import "todo/internal/domain/model"

type TaskUsecase interface {
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task) error
	DeleteTask(id string) error
	GetTask(id string) (*model.Task, error)
	ListTask() ([]*model.Task, error)
}
