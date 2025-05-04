package usecase

import "todo/internal/domain/model"

type TaskUsecase interface {
	CreateTask(task *model.Task) (*model.Task, error)
	UpdateTask(task *model.Task) (*model.Task, error)
	DeleteTask(id string) error
	GetTask(id string) (*model.Task, error)
	ListTasks() ([]*model.Task, error)
	SetTaskCompletion(task *model.Task) (*model.Task, error)
}
