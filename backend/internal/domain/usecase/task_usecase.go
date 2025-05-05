package usecase

import (
	"todo/internal/domain/model"
)

type TaskUsecase interface {
	CreateTask(task *model.Task) (*model.Task, error)
	UpdateTask(task *model.Task) (*model.Task, error)
	DeleteTask(id string) error
	GetTask(id string) (*model.Task, error)
	ListTasksWithFilter(filter *model.TaskFilter) ([]*model.Task, int, error)
	SetTaskCompletion(task *model.Task) (*model.Task, error)
	UpdateOverdueTasks() error
}
